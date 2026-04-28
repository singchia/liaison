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

type View = "main" | "settings";

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

// Heroicons-style cog: a circular hub with 8 trapezoidal teeth and an
// inner ring. Reads as a gear at 14px instead of looking like a sun.
function GearIcon() {
  return (
    <svg
      viewBox="0 0 24 24"
      fill="currentColor"
      aria-hidden="true"
    >
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M11.078 2.25c-.917 0-1.699.663-1.85 1.567L9.05 4.889c-.02.12-.115.26-.297.348a7.493 7.493 0 0 0-.986.57c-.166.115-.334.126-.45.083L6.3 5.508a1.875 1.875 0 0 0-2.282.819l-.922 1.597a1.875 1.875 0 0 0 .432 2.385l.84.692c.095.078.17.229.154.43a7.598 7.598 0 0 0 0 1.139c.015.2-.059.352-.153.43l-.841.692a1.875 1.875 0 0 0-.432 2.385l.922 1.597a1.875 1.875 0 0 0 2.282.818l1.019-.382c.115-.043.283-.031.45.082.312.214.641.405.985.57.182.088.277.228.297.35l.179 1.07c.151.905.933 1.568 1.85 1.568h1.844c.916 0 1.699-.663 1.85-1.567l.178-1.072c.02-.12.114-.26.297-.349.344-.165.673-.356.985-.57.167-.114.335-.125.45-.082l1.02.382a1.875 1.875 0 0 0 2.28-.819l.923-1.597a1.875 1.875 0 0 0-.432-2.385l-.84-.692c-.095-.078-.17-.229-.154-.43a7.55 7.55 0 0 0 0-1.139c-.016-.2.059-.352.153-.43l.84-.692c.708-.582.891-1.59.433-2.385l-.922-1.597a1.875 1.875 0 0 0-2.282-.818l-1.02.382c-.114.043-.282.031-.449-.083a7.49 7.49 0 0 0-.985-.57c-.183-.087-.277-.227-.297-.348L16.772 3.817a1.875 1.875 0 0 0-1.85-1.567h-3.844ZM12 15.75a3.75 3.75 0 1 0 0-7.5 3.75 3.75 0 0 0 0 7.5Z"
      />
    </svg>
  );
}

// Left-pointing chevron for the settings page back button.
function BackArrowIcon() {
  return (
    <svg
      viewBox="0 0 16 16"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.6"
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
    >
      <path d="M10 12.5 5.5 8 10 3.5" />
    </svg>
  );
}

interface MainProps {
  status: StatusPayload;
  busy: boolean;
  onCommand: (cmd: string) => void;
  onOpenSettings: () => void;
  cliBannerDismissed: boolean;
  onDismissCliBanner: () => void;
}

function MainView({
  status,
  busy,
  onCommand,
  onOpenSettings,
  cliBannerDismissed,
  onDismissCliBanner,
}: MainProps) {
  const { tray, logged_in, cli_hits } = status;
  const errMsg = tray.kind === "Error" ? tray.message : null;
  const showCliBanner = cli_hits.length > 0 && !cliBannerDismissed;

  return (
    <main className="popup">
      <header className="hdr" data-tauri-drag-region>
        <img className="hdr__logo" src="/liaison-mark.svg" alt="Liaison" data-tauri-drag-region />
        <span className="hdr__title" data-tauri-drag-region>Liaison</span>
        <span className="hdr__spacer" data-tauri-drag-region />
        <span className={DOT_CLASS[tray.kind]} data-tauri-drag-region />
        <span className="label" data-tauri-drag-region>{LABELS[tray.kind]}</span>
      </header>

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
            onClick={onDismissCliBanner}
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
            onClick={() => onCommand("cmd_login")}
          >
            登录 Liaison
          </button>
        )}

        {logged_in && tray.kind === "Paused" && (
          <button
            className="btn primary"
            disabled={busy}
            onClick={() => onCommand("cmd_resume")}
          >
            恢复连接
          </button>
        )}

        {logged_in && (tray.kind === "Online" || tray.kind === "Connecting") && (
          <button
            className="btn"
            disabled={busy}
            onClick={() => onCommand("cmd_pause")}
          >
            暂停连接
          </button>
        )}

        {logged_in && tray.kind === "Error" && (
          <button
            className="btn primary"
            disabled={busy}
            onClick={() => onCommand("cmd_resume")}
          >
            重连
          </button>
        )}

        {logged_in && (
          <button
            className="btn"
            disabled={busy}
            onClick={() => onCommand("cmd_open_dashboard")}
          >
            打开 Dashboard
          </button>
        )}

        {logged_in && (
          <button
            className="btn ghost"
            disabled={busy}
            onClick={() => onCommand("cmd_logout")}
          >
            退出登录
          </button>
        )}
      </section>

      <footer className="footer">
        <button
          className="footer__gear"
          onClick={onOpenSettings}
          disabled={busy}
          aria-label="设置"
          title="设置"
        >
          <GearIcon />
        </button>
      </footer>
    </main>
  );
}

interface SettingsProps {
  baseUrl: string;
  busy: boolean;
  onClose: () => void;
  onSave: (newUrl: string) => Promise<void>;
}

function SettingsView({ baseUrl, busy, onClose, onSave }: SettingsProps) {
  const [draft, setDraft] = useState(baseUrl);
  const [error, setError] = useState<string | null>(null);

  async function handleSave() {
    const trimmed = draft.trim();
    if (!trimmed) {
      setError("地址不能为空");
      return;
    }
    setError(null);
    try {
      await onSave(trimmed);
    } catch (err) {
      setError(String(err));
    }
  }

  return (
    <main className="popup">
      <header className="hdr" data-tauri-drag-region>
        <button
          className="hdr__back"
          onClick={onClose}
          disabled={busy}
          aria-label="返回"
          title="返回"
        >
          <BackArrowIcon />
        </button>
        <span className="hdr__title" data-tauri-drag-region>服务器配置</span>
        <span className="hdr__spacer" data-tauri-drag-region />
      </header>

      <div className="settings">
        <label className="settings__label" htmlFor="server-url">
          Liaison 服务器地址
        </label>
        <input
          id="server-url"
          className="settings__input"
          type="text"
          value={draft}
          placeholder="https://liaison.example.com"
          onChange={(e) => setDraft(e.target.value)}
          disabled={busy}
          autoFocus
        />
        <div className="settings__hint">
          公网用户保持默认 https://liaison.cloud。私有化部署请填写你的部署地址，
          以 http:// 或 https:// 开头。
        </div>
        {error && <div className="err settings__err">{error}</div>}
      </div>

      <section className="actions">
        <button
          className="btn primary"
          disabled={busy}
          onClick={handleSave}
        >
          保存并重新登录
        </button>
        <button
          className="btn ghost"
          disabled={busy}
          onClick={onClose}
        >
          取消
        </button>
      </section>
    </main>
  );
}

function App() {
  const [status, setStatus] = useState<StatusPayload | null>(null);
  const [busy, setBusy] = useState(false);
  const [cliBannerDismissed, setCliBannerDismissed] = useState(false);
  const [view, setView] = useState<View>("main");

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

  async function runCommand(cmd: string) {
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

  async function saveServer(newUrl: string) {
    setBusy(true);
    try {
      await invoke("cmd_set_server", { newBaseUrl: newUrl });
      setView("main");
      await refresh();
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

  if (view === "settings") {
    return (
      <SettingsView
        baseUrl={status.base_url}
        busy={busy}
        onClose={() => setView("main")}
        onSave={saveServer}
      />
    );
  }

  return (
    <MainView
      status={status}
      busy={busy}
      onCommand={runCommand}
      onOpenSettings={() => setView("settings")}
      cliBannerDismissed={cliBannerDismissed}
      onDismissCliBanner={() => setCliBannerDismissed(true)}
    />
  );
}

export default App;
