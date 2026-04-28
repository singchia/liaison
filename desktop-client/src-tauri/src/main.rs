// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

fn main() {
    // Headless mode invoked by the NSIS uninstaller (see
    // installer-hooks.nsh). Wipes the per-user PAT from the OS
    // credential store before the installer deletes program files,
    // so a clean uninstall doesn't leak a stale credential entry.
    let args: Vec<String> = std::env::args().collect();
    if args.iter().any(|a| a == "--cleanup-credentials") {
        desktop_client_lib::cleanup_credentials();
        return;
    }
    desktop_client_lib::run()
}
