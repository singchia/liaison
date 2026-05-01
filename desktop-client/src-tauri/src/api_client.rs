// HTTP client for liaison.cloud /api/v1/* using PAT bearer auth.
// Server-side contracts:
//   GET  /api/v1/iam/profile  -> Profile
//   GET  /api/v1/edges        -> list edges (paged)
//   POST /api/v1/edges        -> create edge, returns plaintext access/secret
//
// Error envelope from the server:
//   { "code": 401, "message": "UNAUTHORIZED", "details": "..." }
// We map 401 -> ApiError::Unauthorized so the GUI can wipe the keychain
// PAT and prompt re-login without showing a generic HTTP error.

#![allow(dead_code)]

use reqwest::{header, StatusCode};
use serde::{Deserialize, Serialize};
use thiserror::Error;

const DEFAULT_TIMEOUT_SECS: u64 = 15;
const USER_AGENT: &str = concat!("liaison-desktop/", env!("CARGO_PKG_VERSION"));

#[derive(Debug, Error)]
pub enum ApiError {
    #[error("http error: {0}")]
    Http(#[from] reqwest::Error),
    #[error("invalid base url: {0}")]
    InvalidUrl(String),
    #[error("unauthorized — PAT is invalid or expired")]
    Unauthorized,
    #[error("forbidden: {0}")]
    Forbidden(String),
    #[error("server error {status}: {message}")]
    Server { status: u16, message: String },
    #[error("malformed response: {0}")]
    Malformed(String),
}

// Manager APIs are protobuf-generated; protojson encodes int64 / uint64 as
// JSON strings (per the protobuf JSON spec) to dodge JS precision loss.
// Accept both shapes so the client survives whichever encoding the manager
// happens to emit.
fn de_u64_flex<'de, D>(d: D) -> Result<u64, D::Error>
where
    D: serde::Deserializer<'de>,
{
    use serde::de::{self, Deserialize};
    #[derive(Deserialize)]
    #[serde(untagged)]
    enum Either {
        Num(u64),
        Str(String),
    }
    match Either::deserialize(d)? {
        Either::Num(n) => Ok(n),
        Either::Str(s) => s.parse::<u64>().map_err(de::Error::custom),
    }
}

#[derive(Debug, Deserialize, Clone)]
pub struct Profile {
    #[serde(deserialize_with = "de_u64_flex")]
    pub id: u64,
    #[serde(default)]
    pub email: String,
    #[serde(default)]
    pub created_at: String,
    #[serde(default)]
    pub last_login: String,
    #[serde(default)]
    pub login_ip: String,
}

#[derive(Debug, Deserialize, Clone)]
pub struct Edge {
    #[serde(deserialize_with = "de_u64_flex")]
    pub id: u64,
    pub name: String,
    #[serde(default)]
    pub description: String,
    #[serde(default)]
    pub status: i32,
    #[serde(default)]
    pub online: i32,
    #[serde(default)]
    pub created_at: String,
    #[serde(default)]
    pub updated_at: String,
    #[serde(default)]
    pub application_count: i32,
}

#[derive(Debug, Clone)]
pub struct EdgeKeys {
    pub access_key: String,
    pub secret_key: String,
    pub install_command: String,
}

#[derive(Debug, Deserialize)]
struct Envelope<T> {
    code: i32,
    #[serde(default)]
    message: String,
    #[serde(default)]
    details: String,
    data: Option<T>,
}

#[derive(Debug, Deserialize)]
struct ListEdgesData {
    #[serde(default)]
    total: i32,
    #[serde(default)]
    edges: Vec<Edge>,
}

#[derive(Debug, Deserialize)]
struct CreateEdgeData {
    access_key: String,
    secret_key: String,
    #[serde(default)]
    command: String,
}

#[derive(Debug, Serialize)]
struct CreateEdgeRequest<'a> {
    name: &'a str,
    description: &'a str,
}

pub struct ApiClient {
    base_url: String,
    client: reqwest::Client,
}

impl ApiClient {
    pub fn new(base_url: impl Into<String>, pat: impl AsRef<str>) -> Result<Self, ApiError> {
        let mut headers = header::HeaderMap::new();
        let auth = format!("Bearer {}", pat.as_ref());
        let mut auth_value = header::HeaderValue::from_str(&auth)
            .map_err(|e| ApiError::InvalidUrl(format!("invalid PAT for header: {e}")))?;
        auth_value.set_sensitive(true);
        headers.insert(header::AUTHORIZATION, auth_value);
        headers.insert(
            header::ACCEPT,
            header::HeaderValue::from_static("application/json"),
        );

        let base_url_owned = base_url.into().trim_end_matches('/').to_string();

        let mut builder = reqwest::Client::builder()
            .default_headers(headers)
            .user_agent(USER_AGENT)
            .timeout(std::time::Duration::from_secs(DEFAULT_TIMEOUT_SECS));

        // Trust validation strategy:
        //   - SaaS host (liaison.cloud): always validate certs.
        //   - Anything else (private / self-hosted): default to skip,
        //     because private deployments very commonly run with self-
        //     signed certs and otherwise cmd_login partially succeeds
        //     (PAT saved to keychain, then create_edge fails on the
        //     TLS handshake) — leaving the user stuck in "logged in
        //     but no edge.yaml" state.
        //   - LIAISON_INSECURE_TLS env var still forces skip on any host.
        let force_skip = std::env::var("LIAISON_INSECURE_TLS")
            .ok()
            .filter(|s| !s.is_empty() && s != "0")
            .is_some();
        let host_is_saas = url::Url::parse(&base_url_owned)
            .ok()
            .and_then(|u| u.host_str().map(str::to_string))
            .map(|h| h == "liaison.cloud")
            .unwrap_or(false);
        if force_skip || !host_is_saas {
            builder = builder.danger_accept_invalid_certs(true);
        }
        let client = builder.build()?;

        Ok(Self {
            base_url: base_url_owned,
            client,
        })
    }

    pub async fn get_profile(&self) -> Result<Profile, ApiError> {
        let url = format!("{}/api/v1/iam/profile", self.base_url);
        let resp = self.client.get(&url).send().await?;
        unwrap_envelope::<Profile>(resp).await
    }

    pub async fn list_edges(&self) -> Result<Vec<Edge>, ApiError> {
        let url = format!("{}/api/v1/edges", self.base_url);
        let resp = self
            .client
            .get(&url)
            .query(&[("page", "1"), ("page_size", "100")])
            .send()
            .await?;
        let data = unwrap_envelope::<ListEdgesData>(resp).await?;
        Ok(data.edges)
    }

    pub async fn create_edge(
        &self,
        name: &str,
        description: &str,
    ) -> Result<EdgeKeys, ApiError> {
        let url = format!("{}/api/v1/edges", self.base_url);
        let body = CreateEdgeRequest { name, description };
        let resp = self.client.post(&url).json(&body).send().await?;
        let data = unwrap_envelope::<CreateEdgeData>(resp).await?;
        Ok(EdgeKeys {
            access_key: data.access_key,
            secret_key: data.secret_key,
            install_command: data.command,
        })
    }
}

async fn unwrap_envelope<T: for<'de> Deserialize<'de>>(
    resp: reqwest::Response,
) -> Result<T, ApiError> {
    let status = resp.status();
    if status == StatusCode::UNAUTHORIZED {
        return Err(ApiError::Unauthorized);
    }

    let bytes = resp.bytes().await?;
    if status == StatusCode::FORBIDDEN {
        let msg = error_message_from_body(&bytes).unwrap_or_else(|| "forbidden".to_string());
        return Err(ApiError::Forbidden(msg));
    }
    if !status.is_success() {
        let msg = error_message_from_body(&bytes)
            .unwrap_or_else(|| format!("HTTP {}", status.as_u16()));
        return Err(ApiError::Server {
            status: status.as_u16(),
            message: msg,
        });
    }

    // The server uses a {code, message, data} envelope even on 200 responses.
    // Treat code != 200 inside a 2xx response as a malformed server reply.
    let env: Envelope<T> = serde_json::from_slice(&bytes)
        .map_err(|e| ApiError::Malformed(format!("decode envelope: {e}")))?;
    if env.code != 0 && env.code != 200 {
        let msg = if env.details.is_empty() {
            env.message
        } else {
            format!("{}: {}", env.message, env.details)
        };
        return Err(ApiError::Server {
            status: status.as_u16(),
            message: msg,
        });
    }
    env.data
        .ok_or_else(|| ApiError::Malformed("missing data field".into()))
}

fn error_message_from_body(bytes: &[u8]) -> Option<String> {
    let v: serde_json::Value = serde_json::from_slice(bytes).ok()?;
    let msg = v.get("message").and_then(|m| m.as_str()).unwrap_or("");
    let details = v.get("details").and_then(|d| d.as_str()).unwrap_or("");
    match (msg.is_empty(), details.is_empty()) {
        (true, true) => None,
        (false, true) => Some(msg.to_string()),
        (true, false) => Some(details.to_string()),
        (false, false) => Some(format!("{msg}: {details}")),
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use wiremock::matchers::{header, method, path, query_param};
    use wiremock::{Mock, MockServer, ResponseTemplate};

    async fn fixture_server() -> MockServer {
        MockServer::start().await
    }

    #[tokio::test]
    async fn get_profile_decodes_envelope() {
        let server = fixture_server().await;
        Mock::given(method("GET"))
            .and(path("/api/v1/iam/profile"))
            .and(header("Authorization", "Bearer liaison_pat_test"))
            .respond_with(ResponseTemplate::new(200).set_body_json(serde_json::json!({
                "code": 200,
                "message": "success",
                "data": {
                    "id": 42,
                    "email": "user@example.com",
                    "created_at": "2026-01-01T00:00:00Z",
                    "last_login": "2026-04-25T10:00:00Z",
                    "login_ip": "10.0.0.1"
                }
            })))
            .mount(&server)
            .await;

        let client = ApiClient::new(server.uri(), "liaison_pat_test").unwrap();
        let p = client.get_profile().await.unwrap();
        assert_eq!(p.id, 42);
        assert_eq!(p.email, "user@example.com");
    }

    #[tokio::test]
    async fn list_edges_returns_array() {
        let server = fixture_server().await;
        Mock::given(method("GET"))
            .and(path("/api/v1/edges"))
            .and(query_param("page", "1"))
            .and(query_param("page_size", "100"))
            .respond_with(ResponseTemplate::new(200).set_body_json(serde_json::json!({
                "code": 200,
                "message": "success",
                "data": {
                    "total": 1,
                    "edges": [{
                        "id": 7,
                        "name": "demo",
                        "description": "",
                        "status": 1,
                        "online": 2,
                        "created_at": "2026-04-25T10:00:00Z",
                        "updated_at": "2026-04-25T10:00:00Z",
                        "application_count": 0
                    }]
                }
            })))
            .mount(&server)
            .await;

        let client = ApiClient::new(server.uri(), "pat").unwrap();
        let edges = client.list_edges().await.unwrap();
        assert_eq!(edges.len(), 1);
        assert_eq!(edges[0].id, 7);
        assert_eq!(edges[0].name, "demo");
    }

    #[tokio::test]
    async fn create_edge_returns_keys() {
        let server = fixture_server().await;
        Mock::given(method("POST"))
            .and(path("/api/v1/edges"))
            .respond_with(ResponseTemplate::new(200).set_body_json(serde_json::json!({
                "code": 200,
                "message": "success",
                "data": {
                    "access_key": "ak_xxx",
                    "secret_key": "sk_yyy",
                    "command": "curl -fsSL ..."
                }
            })))
            .mount(&server)
            .await;

        let client = ApiClient::new(server.uri(), "pat").unwrap();
        let keys = client
            .create_edge("Liaison Desktop", "menubar app")
            .await
            .unwrap();
        assert_eq!(keys.access_key, "ak_xxx");
        assert_eq!(keys.secret_key, "sk_yyy");
        assert!(keys.install_command.starts_with("curl"));
    }

    #[tokio::test]
    async fn unauthorized_maps_to_typed_error() {
        let server = fixture_server().await;
        Mock::given(method("GET"))
            .and(path("/api/v1/iam/profile"))
            .respond_with(ResponseTemplate::new(401).set_body_json(serde_json::json!({
                "code": 401,
                "message": "UNAUTHORIZED",
                "details": "Token validation failed"
            })))
            .mount(&server)
            .await;

        let client = ApiClient::new(server.uri(), "stale").unwrap();
        let err = client.get_profile().await.unwrap_err();
        assert!(matches!(err, ApiError::Unauthorized));
    }

    #[tokio::test]
    async fn forbidden_carries_message() {
        let server = fixture_server().await;
        Mock::given(method("GET"))
            .and(path("/api/v1/iam/profile"))
            .respond_with(ResponseTemplate::new(403).set_body_json(serde_json::json!({
                "code": 403,
                "message": "ACCOUNT_DELETED",
                "details": "User already deleted"
            })))
            .mount(&server)
            .await;

        let client = ApiClient::new(server.uri(), "pat").unwrap();
        let err = client.get_profile().await.unwrap_err();
        match err {
            ApiError::Forbidden(msg) => assert!(msg.contains("ACCOUNT_DELETED")),
            other => panic!("expected Forbidden, got {other:?}"),
        }
    }

    #[tokio::test]
    async fn list_edges_accepts_string_encoded_ids() {
        // Real production manager (protobuf-generated) returns uint64
        // ids as strings. The list_edges call was failing in prod with
        // "invalid type: string \"100152\", expected u64" until the
        // flexible deserializer was added — this guards that regression.
        let server = fixture_server().await;
        Mock::given(method("GET"))
            .and(path("/api/v1/edges"))
            .respond_with(ResponseTemplate::new(200).set_body_json(serde_json::json!({
                "code": 200,
                "message": "success",
                "data": {
                    "total": 1,
                    "edges": [{
                        "id": "100152",
                        "name": "Liaison Desktop",
                        "online": 1
                    }]
                }
            })))
            .mount(&server)
            .await;
        let client = ApiClient::new(server.uri(), "pat").unwrap();
        let edges = client.list_edges().await.unwrap();
        assert_eq!(edges[0].id, 100152);
        assert_eq!(edges[0].online, 1);
    }

    #[tokio::test]
    async fn malformed_data_field_errors() {
        let server = fixture_server().await;
        Mock::given(method("GET"))
            .and(path("/api/v1/iam/profile"))
            .respond_with(ResponseTemplate::new(200).set_body_json(serde_json::json!({
                "code": 200,
                "message": "success"
                // no data field
            })))
            .mount(&server)
            .await;

        let client = ApiClient::new(server.uri(), "pat").unwrap();
        let err = client.get_profile().await.unwrap_err();
        assert!(matches!(err, ApiError::Malformed(_)));
    }
}
