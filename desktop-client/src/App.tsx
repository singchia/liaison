import { useEffect, useState } from "react";
import { invoke } from "@tauri-apps/api/core";
import { listen } from "@tauri-apps/api/event";
import "./App.css";

type TrayState =
  | { kind: "LoggedOut" }
  | { kind: "Paused" }
  | { kind: "Connecting" }
  | { kind: "Online" }
  | { kind: "Error"; message: string };

interface CliHit {
  kind: "binary" | "config" | "autostart";
  detail: string;
}

interface StatusPayload {
  tray: TrayState;
  logged_in: boolean;
  cli_hits: CliHit[];
  base_url: string;
}

const LABELS: Record<TrayState["kind"], string> = {
  LoggedOut: "未登录",
  Paused: "已暂停",
  Connecting: "连接中",
  Online: "已连接",
  Error: "错误",
};

const DOT_CLASS: Record<TrayState["kind"], string> = {
  LoggedOut: "dot dot--grey",
  Paused: "dot dot--grey",
  Connecting: "dot dot--amber",
  Online: "dot dot--green",
  Error: "dot dot--red",
};

function hostFromUrl(url: string): string {
  try {
    return new URL(url).host;
  } catch {
    return url;
  }
}

function App() {
  const [status, setStatus] = useState<StatusPayload | null>(null);
  const [busy, setBusy] = useState(false);
  const [cliBannerDismissed, setCliBannerDismissed] = useState(false);
  const [editingServer, setEditingServer] = useState(false);
  const [serverDraft, setServerDraft] = useState("");
  const [serverError, setServerError] = useState<string | null>(null);

  async function refresh() {
    try {
      const s = await invoke<StatusPayload>("cmd_get_status");
      setStatus(s);
    } catch (err) {
      console.error("get_status failed", err);
    }
  }

  useEffect(() => {
    refresh();
    const unlistenP = listen<TrayState>("tray_state_changed", (event) => {
      setStatus((prev) =>
        prev
          ? { ...prev, tray: event.payload }
          : { tray: event.payload, logged_in: true, cli_hits: [], base_url: "" }
      );
      // Server-switch path emits LoggedOut; pull a fresh full payload
      // so base_url, logged_in, etc. all repaint together.
      if (event.payload.kind === "LoggedOut") refresh();
    });
    return () => {
      unlistenP.then((fn) => fn()).catch(() => {});
    };
  }, []);

  async function run(cmd: string) {
    setBusy(true);
    try {
      await invoke(cmd);
    } catch (err) {
      console.error(`${cmd} failed`, err);
    } finally {
      setBusy(false);
      refresh();
    }
  }

  function startEditServer(current: string) {
    setServerDraft(current);
    setServerError(null);
    setEditingServer(true);
  }

  async function saveServer() {
    const trimmed = serverDraft.trim();
    if (!trimmed) {
      setServerError("地址不能为空");
      return;
    }
    setBusy(true);
    setServerError(null);
    try {
      await invoke("cmd_set_server", { newBaseUrl: trimmed });
      setEditingServer(false);
      refresh();
    } catch (err) {
      setServerError(String(err));
    } finally {
      setBusy(false);
    }
  }

  if (!status) {
    return (
      <main className="popup">
        <div className="loading">…</div>
      </main>
    );
  }

  const { tray, logged_in, cli_hits, base_url } = status;
  const errMsg = tray.kind === "Error" ? tray.message : null;
  const showCliBanner = cli_hits.length > 0 && !cliBannerDismissed;
  const host = hostFromUrl(base_url);

  return (
    <main className="popup">
      <header className="hdr" data-tauri-drag-region>
        <img className="hdr__logo" src="/liaison-mark.svg" alt="Liaison" data-tauri-drag-region />
        <span className="hdr__title" data-tauri-drag-region>Liaison</span>
        <span className="hdr__spacer" data-tauri-drag-region />
        <span className={DOT_CLASS[tray.kind]} data-tauri-drag-region />
        <span className="label" data-tauri-drag-region>{LABELS[tray.kind]}</span>
      </header>

      {editingServer ? (
        <div className="server server--edit">
          <input
            className="server__input"
            type="text"
            value={serverDraft}
            placeholder="https://liaison.example.com"
            onChange={(e) => setServerDraft(e.target.value)}
            disabled={busy}
            autoFocus
          />
          {serverError && <div className="server__err">{serverError}</div>}
          <div className="server__row">
            <button
              className="btn primary server__btn"
              disabled={busy}
              onClick={saveServer}
            >
              保存并重新登录
            </button>
            <button
              className="btn ghost server__btn"
              disabled={busy}
              onClick={() => setEditingServer(false)}
            >
              取消
            </button>
          </div>
        </div>
      ) : (
        <div className="server">
          <span className="server__label">服务器</span>
          <span className="server__host" title={base_url}>{host || "—"}</span>
          <button
            className="server__change"
            disabled={busy}
            onClick={() => startEditServer(base_url)}
          >
            更改
          </button>
        </div>
      )}

      {showCliBanner && (
        <div className="banner">
          <div className="banner__title">
            ⚠ 检测到本机已有 liaison-edge CLI 安装
          </div>
          <div className="banner__body">
            继续使用 Liaison Desktop 会创建一个独立的连接器，与 CLI 并行运行。
            建议先卸载 CLI 版本：
          </div>
          <ul className="banner__list">
            {cli_hits.slice(0, 4).map((h, i) => (
              <li key={i}>
                <span className="banner__kind">{h.kind}</span>
                <span className="banner__detail">{h.detail}</span>
              </li>
            ))}
            {cli_hits.length > 4 && (
              <li className="banner__more">
                + {cli_hits.length - 4} 项…
              </li>
            )}
          </ul>
          <button
            className="btn ghost banner__dismiss"
            onClick={() => setCliBannerDismissed(true)}
          >
            知道了，继续
          </button>
        </div>
      )}

      {errMsg && <div className="err">{errMsg}</div>}

      <section className="actions">
        {!logged_in && (
          <button
            className="btn primary"
            disabled={busy}
            onClick={() => run("cmd_login")}
          >
            登录 Liaison
          </button>
        )}

        {logged_in && tray.kind === "Paused" && (
          <button
            className="btn primary"
            disabled={busy}
            onClick={() => run("cmd_resume")}
          >
            恢复连接
          </button>
        )}

        {logged_in && (tray.kind === "Online" || tray.kind === "Connecting") && (
          <button
            className="btn"
            disabled={busy}
            onClick={() => run("cmd_pause")}
          >
            暂停连接
          </button>
        )}

        {logged_in && tray.kind === "Error" && (
          <button
            className="btn primary"
            disabled={busy}
            onClick={() => run("cmd_resume")}
          >
            重连
          </button>
        )}

        {logged_in && (
          <button
            className="btn"
            disabled={busy}
            onClick={() => run("cmd_open_dashboard")}
          >
            打开 Dashboard
          </button>
        )}

        {logged_in && (
          <button
            className="btn ghost"
            disabled={busy}
            onClick={() => run("cmd_logout")}
          >
            退出登录
          </button>
        )}
      </section>
    </main>
  );
}

export default App;
