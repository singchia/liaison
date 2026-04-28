// Device-code-style login: spin up a localhost HTTP server, point the
// browser at /dashboard/cli-auth, receive the PAT back via the callback
// query string, store it in the OS keychain.
//
// Server-side contract is documented in:
//   web/src/pages/CliAuth/index.tsx           (callback regex, params)
//   pkg/liaison/manager/web/iam_tokens_http.go (POST /api/v1/iam/tokens)
//   pkg/liaison/iam/pat.go                     (token format, prefix)

#![allow(dead_code)]

use std::time::Duration;

use rand::Rng;
use thiserror::Error;
use tokio::sync::oneshot;

const KEYCHAIN_SERVICE: &str = "cloud.liaison.desktop";
const KEYCHAIN_USER_LEGACY: &str = "pat";
const STATE_LEN: usize = 32;

/// Keychain account name namespaced by server host so a user with both a
/// SaaS and a self-hosted deployment doesn't share the PAT slot.
fn keychain_user(base_url: &str) -> String {
    match url::Url::parse(base_url).ok().and_then(|u| u.host_str().map(str::to_string)) {
        Some(host) if !host.is_empty() => format!("pat@{host}"),
        _ => KEYCHAIN_USER_LEGACY.to_string(),
    }
}
const DEFAULT_TOKEN_NAME: &str = "Liaison Desktop";
const CALLBACK_PATH: &str = "/callback";
const DEFAULT_LOGIN_TIMEOUT_SECS: u64 = 300;

#[derive(Debug, Error)]
pub enum AuthError {
    #[error("io error: {0}")]
    Io(#[from] std::io::Error),
    #[error("keychain error: {0}")]
    Keyring(#[from] keyring::Error),
    #[error("invalid base url: {0}")]
    InvalidUrl(String),
    #[error("login server error: {0}")]
    Server(String),
    #[error("login timed out")]
    Timeout,
    #[error("login was cancelled")]
    Cancelled,
    #[error("state mismatch — possible CSRF")]
    StateMismatch,
    #[error("login returned no token")]
    NoToken,
    #[error("login error from server: {0}")]
    RemoteError(String),
}

pub struct PendingLogin {
    pub auth_url: String,
    rx: oneshot::Receiver<Result<String, AuthError>>,
}

pub fn start_login(base_url: &str, token_name: Option<&str>) -> Result<PendingLogin, AuthError> {
    let server = tiny_http::Server::http("127.0.0.1:0")
        .map_err(|e| AuthError::Server(e.to_string()))?;
    let port = server
        .server_addr()
        .to_ip()
        .ok_or_else(|| AuthError::Server("listener is not TCP".into()))?
        .port();

    let state = generate_state();
    let name = token_name.unwrap_or(DEFAULT_TOKEN_NAME);
    let auth_url = build_auth_url(base_url, port, &state, name)?;
    let expected_state = state;

    let (tx, rx) = oneshot::channel::<Result<String, AuthError>>();

    std::thread::spawn(move || {
        for request in server.incoming_requests() {
            let url = request.url().to_string();
            if !url.starts_with(CALLBACK_PATH) {
                let _ = request.respond(tiny_http::Response::empty(404));
                continue;
            }
            let result = parse_callback(&url, &expected_state);
            let body = match &result {
                Ok(_) => CALLBACK_OK_HTML,
                Err(AuthError::Cancelled) => CALLBACK_CANCELLED_HTML,
                Err(_) => CALLBACK_ERROR_HTML,
            };
            let resp = tiny_http::Response::from_string(body).with_header(
                tiny_http::Header::from_bytes(
                    &b"Content-Type"[..],
                    &b"text/html; charset=utf-8"[..],
                )
                .expect("static header"),
            );
            let _ = request.respond(resp);
            let _ = tx.send(result);
            return;
        }
    });

    Ok(PendingLogin { auth_url, rx })
}

impl PendingLogin {
    pub async fn await_token(self, timeout: Duration) -> Result<String, AuthError> {
        match tokio::time::timeout(timeout, self.rx).await {
            Ok(Ok(inner)) => inner,
            Ok(Err(_)) => Err(AuthError::Cancelled),
            Err(_) => Err(AuthError::Timeout),
        }
    }
}

pub async fn login_with_device_code<F>(
    base_url: &str,
    open_browser: F,
) -> Result<String, AuthError>
where
    F: FnOnce(&str) -> Result<(), AuthError>,
{
    let pending = start_login(base_url, None)?;
    open_browser(&pending.auth_url)?;
    let pat = pending
        .await_token(Duration::from_secs(DEFAULT_LOGIN_TIMEOUT_SECS))
        .await?;
    save_pat_to_keychain(base_url, &pat)?;
    Ok(pat)
}

pub fn save_pat_to_keychain(base_url: &str, pat: &str) -> Result<(), AuthError> {
    let entry = keyring::Entry::new(KEYCHAIN_SERVICE, &keychain_user(base_url))?;
    entry.set_password(pat)?;
    Ok(())
}

/// Read the PAT for `base_url`. Falls back to the legacy un-namespaced
/// entry so users upgrading from the single-server build don't get
/// logged out the first time they launch the multi-server build.
pub fn get_pat_from_keychain(base_url: &str) -> Result<String, AuthError> {
    let user = keychain_user(base_url);
    let entry = keyring::Entry::new(KEYCHAIN_SERVICE, &user)?;
    match entry.get_password() {
        Ok(pat) => Ok(pat),
        Err(keyring::Error::NoEntry) if user != KEYCHAIN_USER_LEGACY => {
            let legacy = keyring::Entry::new(KEYCHAIN_SERVICE, KEYCHAIN_USER_LEGACY)?;
            Ok(legacy.get_password()?)
        }
        Err(e) => Err(e.into()),
    }
}

pub fn delete_pat_from_keychain(base_url: &str) -> Result<(), AuthError> {
    let entry = keyring::Entry::new(KEYCHAIN_SERVICE, &keychain_user(base_url))?;
    entry.delete_credential()?;
    Ok(())
}

fn parse_callback(path_with_query: &str, expected_state: &str) -> Result<String, AuthError> {
    let parsed = url::Url::parse(&format!("http://localhost{path_with_query}"))
        .map_err(|e| AuthError::InvalidUrl(e.to_string()))?;
    let mut state = None;
    let mut token = None;
    let mut error_msg = None;
    for (k, v) in parsed.query_pairs() {
        match k.as_ref() {
            "state" => state = Some(v.into_owned()),
            "token" => token = Some(v.into_owned()),
            "error" => error_msg = Some(v.into_owned()),
            _ => {}
        }
    }
    if let Some(err) = error_msg {
        return match err.as_str() {
            "denied" | "cancelled" | "user_denied" => Err(AuthError::Cancelled),
            other => Err(AuthError::RemoteError(other.to_string())),
        };
    }
    let state = state.ok_or(AuthError::NoToken)?;
    if state != expected_state {
        return Err(AuthError::StateMismatch);
    }
    let token = token.ok_or(AuthError::NoToken)?;
    if token.is_empty() {
        return Err(AuthError::NoToken);
    }
    Ok(token)
}

fn build_auth_url(base: &str, port: u16, state: &str, name: &str) -> Result<String, AuthError> {
    let mut url = url::Url::parse(base).map_err(|e| AuthError::InvalidUrl(e.to_string()))?;
    let trimmed = url.path().trim_end_matches('/').to_string();
    url.set_path(&format!("{trimmed}/dashboard/cli-auth"));
    let callback = format!("http://127.0.0.1:{port}/callback");
    url.query_pairs_mut()
        .clear()
        .append_pair("callback", &callback)
        .append_pair("state", state)
        .append_pair("name", name)
        .append_pair("mode", "callback");
    Ok(url.into())
}

fn generate_state() -> String {
    const CHARSET: &[u8] = b"abcdefghijklmnopqrstuvwxyz0123456789";
    let mut rng = rand::thread_rng();
    (0..STATE_LEN)
        .map(|_| CHARSET[rng.gen_range(0..CHARSET.len())] as char)
        .collect()
}

const CALLBACK_OK_HTML: &str = r#"<!doctype html><html><head><meta charset="utf-8"><title>Liaison · 已登录</title>
<style>body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;background:#0c0e1a;color:#e9eef9;display:flex;align-items:center;justify-content:center;height:100vh;margin:0}main{text-align:center}h1{font-size:18px;margin:0 0 8px}p{font-size:14px;color:#9aa3b8;margin:0}</style></head>
<body><main><h1>已登录 Liaison Desktop</h1><p>可以关闭此页面，返回菜单栏。</p></main></body></html>"#;

const CALLBACK_ERROR_HTML: &str = r#"<!doctype html><html><head><meta charset="utf-8"><title>Liaison · 登录失败</title>
<style>body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;background:#0c0e1a;color:#e9eef9;display:flex;align-items:center;justify-content:center;height:100vh;margin:0}main{text-align:center}h1{font-size:18px;margin:0 0 8px}p{font-size:14px;color:#9aa3b8;margin:0}</style></head>
<body><main><h1>登录失败</h1><p>请回到菜单栏重试。</p></main></body></html>"#;

const CALLBACK_CANCELLED_HTML: &str = r#"<!doctype html><html><head><meta charset="utf-8"><title>Liaison · 已取消</title>
<style>body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;background:#0c0e1a;color:#e9eef9;display:flex;align-items:center;justify-content:center;height:100vh;margin:0}main{text-align:center}h1{font-size:18px;margin:0 0 8px}p{font-size:14px;color:#9aa3b8;margin:0}</style></head>
<body><main><h1>已取消</h1><p>未授权 Liaison Desktop。</p></main></body></html>"#;

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn build_auth_url_appends_path_and_params() {
        let url = build_auth_url("https://liaison.cloud", 47823, "abc1234567890123", "Demo")
            .expect("build url");
        assert!(url.starts_with("https://liaison.cloud/dashboard/cli-auth?"));
        assert!(url.contains("callback=http%3A%2F%2F127.0.0.1%3A47823%2Fcallback"));
        assert!(url.contains("state=abc1234567890123"));
        assert!(url.contains("name=Demo"));
        assert!(url.contains("mode=callback"));
    }

    #[test]
    fn build_auth_url_handles_trailing_slash() {
        let url =
            build_auth_url("https://liaison.cloud/", 1, "s", "n").expect("build url");
        assert!(url.starts_with("https://liaison.cloud/dashboard/cli-auth?"));
    }

    #[test]
    fn parse_callback_returns_token_on_success() {
        let token =
            parse_callback("/callback?state=xyz&token=liaison_pat_abc", "xyz").unwrap();
        assert_eq!(token, "liaison_pat_abc");
    }

    #[test]
    fn parse_callback_rejects_state_mismatch() {
        let err = parse_callback("/callback?state=other&token=t", "expected").unwrap_err();
        assert!(matches!(err, AuthError::StateMismatch));
    }

    #[test]
    fn parse_callback_rejects_missing_token() {
        let err = parse_callback("/callback?state=xyz", "xyz").unwrap_err();
        assert!(matches!(err, AuthError::NoToken));
    }

    #[test]
    fn parse_callback_maps_denied_to_cancelled() {
        let err = parse_callback("/callback?error=denied", "any").unwrap_err();
        assert!(matches!(err, AuthError::Cancelled));
    }

    #[test]
    fn parse_callback_keeps_unknown_errors_as_remote() {
        let err = parse_callback("/callback?error=server_blew_up", "any").unwrap_err();
        match err {
            AuthError::RemoteError(msg) => assert_eq!(msg, "server_blew_up"),
            other => panic!("expected RemoteError, got {other:?}"),
        }
    }

    #[test]
    fn generate_state_is_long_enough() {
        let s = generate_state();
        assert_eq!(s.len(), STATE_LEN);
        assert!(s.chars().all(|c| c.is_ascii_alphanumeric()));
    }
}
