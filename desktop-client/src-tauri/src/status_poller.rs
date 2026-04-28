// Periodically pings GET /api/v1/edges and tells the supervisor whether
// our edge is actually connected to the manager. The supervisor only
// knows whether the subprocess is alive — this poller is what
// distinguishes Connecting (process alive but no tunnel) from Online
// (process alive AND manager has marked online=1).

#![allow(dead_code)]

use std::time::Duration;

use crate::api_client::{ApiClient, ApiError, Edge};
use crate::edge_supervisor::SupervisorHandle;

const POLL_INTERVAL_SECS: u64 = 5;
const ONLINE_FLAG: i32 = 1;

pub fn spawn_poller(
    base_url: String,
    pat: String,
    edge_name: String,
    handle: SupervisorHandle,
) {
    tauri::async_runtime::spawn(async move {
        let api = match ApiClient::new(&base_url, &pat) {
            Ok(api) => api,
            Err(e) => {
                crate::debug_log(format!("poller: build api client failed: {e}"));
                return;
            }
        };

        crate::debug_log(format!(
            "poller: started, base_url={base_url}, edge_name={edge_name}"
        ));

        loop {
            match api.list_edges().await {
                Ok(edges) => {
                    let pick = pick_our_edge(&edges, &edge_name);
                    crate::debug_log(format!(
                        "poller: list_edges total={}, picked={:?}",
                        edges.len(),
                        pick.map(|e| (e.id, e.name.clone(), e.online))
                    ));
                    match pick {
                        Some(edge) if edge.online == ONLINE_FLAG => {
                            handle.report_tunnel_online();
                        }
                        Some(_) => {
                            handle.report_tunnel_offline();
                        }
                        None => {
                            // No edges at all on the account. Don't
                            // change state; supervisor will report its
                            // own subprocess-level error if relevant.
                        }
                    }
                }
                Err(ApiError::Unauthorized) => {
                    crate::debug_log("poller: PAT rejected, exiting");
                    return;
                }
                Err(e) => {
                    crate::debug_log(format!("poller: list_edges failed: {e}"));
                }
            }

            tokio::time::sleep(Duration::from_secs(POLL_INTERVAL_SECS)).await;
        }
    });
}

/// Pick which edge in the list represents the local install. We try in
/// order of preference:
///   1. Exact name match (the name we created during cmd_login). Among
///      multiple matches, the one with the highest ID wins so a user
///      who clicked login twice still tracks the freshest one.
///   2. Otherwise, fall back to the most recently created edge by ID.
///      Robust against case / whitespace differences the server might
///      apply, and against the user having renamed the edge in the
///      dashboard.
fn pick_our_edge<'a>(edges: &'a [Edge], edge_name: &str) -> Option<&'a Edge> {
    let by_name = edges
        .iter()
        .filter(|e| e.name == edge_name)
        .max_by_key(|e| e.id);
    if by_name.is_some() {
        return by_name;
    }
    edges.iter().max_by_key(|e| e.id)
}

#[cfg(test)]
mod tests {
    use super::*;

    fn edge(id: u64, name: &str, online: i32) -> Edge {
        Edge {
            id,
            name: name.into(),
            description: String::new(),
            status: 1,
            online,
            created_at: String::new(),
            updated_at: String::new(),
            application_count: 0,
        }
    }

    #[test]
    fn picks_max_id_with_matching_name() {
        let edges = vec![
            edge(1, "Liaison Desktop", 1),
            edge(5, "Liaison Desktop", 1),
            edge(3, "Other", 1),
        ];
        let pick = pick_our_edge(&edges, "Liaison Desktop").unwrap();
        assert_eq!(pick.id, 5);
    }

    #[test]
    fn falls_back_to_max_id_when_no_name_match() {
        let edges = vec![edge(1, "old-cli", 1), edge(7, "renamed-by-user", 1)];
        let pick = pick_our_edge(&edges, "Liaison Desktop").unwrap();
        assert_eq!(pick.id, 7);
    }

    #[test]
    fn returns_none_for_empty() {
        let edges: Vec<Edge> = vec![];
        assert!(pick_our_edge(&edges, "Liaison Desktop").is_none());
    }
}
