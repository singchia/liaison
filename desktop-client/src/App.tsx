import { useEffect, useState } from "react";
import { invoke } from "@tauri-apps/api/core";
import { listen } from "@tauri-apps/api/event";
import "./App.css";
import { Dict, Locale, detectLocale, dict as loadDict } from "./i18n";

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
  locale: Locale | null;
}

type View = "main" | "settings";

function labelFor(t: Dict, kind: TrayState["kind"]): string {
  switch (kind) {
    case "LoggedOut": return t.label_logged_out;
    case "Paused": return t.label_paused;
    case "Connecting": return t.label_connecting;
    case "Online": return t.label_online;
    case "Error": return t.label_error;
  }
}

const DOT_CLASS: Record<TrayState["kind"], string> = {
  LoggedOut: "dot dot--grey",
  Paused: "dot dot--grey",
  Connecting: "dot dot--amber",
  Online: "dot dot--green",
  Error: "dot dot--red",
};

// Heroicons-style cog: a circular hub with 8 trapezoidal teeth and an
// inner ring. Reads as a gear at 16px instead of looking like a sun.
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
  t: Dict;
  busy: boolean;
  onCommand: (cmd: string) => void;
  onOpenSettings: () => void;
  cliBannerDismissed: boolean;
  onDismissCliBanner: () => void;
}

function MainView({
  status,
  t,
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
        <span className="label" data-tauri-drag-region>{labelFor(t, tray.kind)}</span>
      </header>

      {showCliBanner && (
        <div className="banner">
          <div className="banner__title">{t.cli_banner_title}</div>
          <div className="banner__body">{t.cli_banner_body}</div>
          <ul className="banner__list">
            {cli_hits.slice(0, 4).map((h, i) => (
              <li key={i}>
                <span className="banner__kind">{h.kind}</span>
                <span className="banner__detail">{h.detail}</span>
              </li>
            ))}
            {cli_hits.length > 4 && (
              <li className="banner__more">{t.cli_banner_more(cli_hits.length - 4)}</li>
            )}
          </ul>
          <button className="btn ghost banner__dismiss" onClick={onDismissCliBanner}>
            {t.cli_banner_dismiss}
          </button>
        </div>
      )}

      {errMsg && <div className="err">{errMsg}</div>}

      <section className="actions">
        {!logged_in && (
          <button className="btn primary" disabled={busy} onClick={() => onCommand("cmd_login")}>
            {t.btn_login}
          </button>
        )}

        {logged_in && tray.kind === "Paused" && (
          <button className="btn primary" disabled={busy} onClick={() => onCommand("cmd_resume")}>
            {t.btn_resume}
          </button>
        )}

        {logged_in && (tray.kind === "Online" || tray.kind === "Connecting") && (
          <button className="btn" disabled={busy} onClick={() => onCommand("cmd_pause")}>
            {t.btn_pause}
          </button>
        )}

        {logged_in && tray.kind === "Error" && (
          <button className="btn primary" disabled={busy} onClick={() => onCommand("cmd_resume")}>
            {t.btn_reconnect}
          </button>
        )}

        {logged_in && (
          <button className="btn" disabled={busy} onClick={() => onCommand("cmd_open_dashboard")}>
            {t.btn_dashboard}
          </button>
        )}

        {logged_in && (
          <button className="btn ghost" disabled={busy} onClick={() => onCommand("cmd_logout")}>
            {t.btn_logout}
          </button>
        )}
      </section>

      <footer className="footer">
        <button
          className="footer__gear"
          onClick={onOpenSettings}
          disabled={busy}
          aria-label={t.settings_open_aria}
          title={t.settings_open_aria}
        >
          <GearIcon />
        </button>
      </footer>
    </main>
  );
}

interface SettingsProps {
  baseUrl: string;
  locale: Locale;
  t: Dict;
  busy: boolean;
  onClose: () => void;
  onSave: (newUrl: string) => Promise<void>;
  onLocaleChange: (l: Locale) => void;
}

function SettingsView({
  baseUrl,
  locale,
  t,
  busy,
  onClose,
  onSave,
  onLocaleChange,
}: SettingsProps) {
  const [draft, setDraft] = useState(baseUrl);
  const [error, setError] = useState<string | null>(null);

  async function handleSave() {
    const trimmed = draft.trim();
    if (!trimmed) {
      setError(t.settings_url_empty);
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
          aria-label={t.settings_back_aria}
          title={t.settings_back_aria}
        >
          <BackArrowIcon />
        </button>
        <span className="hdr__title" data-tauri-drag-region>{t.settings_title}</span>
        <span className="hdr__spacer" data-tauri-drag-region />
      </header>

      <div className="settings">
        <label className="settings__label" htmlFor="server-url">{t.settings_url_label}</label>
        <input
          id="server-url"
          className="settings__input"
          type="text"
          value={draft}
          placeholder={t.settings_url_placeholder}
          onChange={(e) => setDraft(e.target.value)}
          disabled={busy}
          autoFocus
        />
        <div className="settings__hint">{t.settings_url_hint}</div>
        {error && <div className="err settings__err">{error}</div>}

        <label className="settings__label settings__label--gap" htmlFor="locale-select">
          {t.settings_locale_label}
        </label>
        <select
          id="locale-select"
          className="settings__select"
          value={locale}
          onChange={(e) => onLocaleChange(e.target.value as Locale)}
          disabled={busy}
        >
          <option value="en">{t.locale_en}</option>
          <option value="zh">{t.locale_zh}</option>
        </select>
      </div>

      <section className="actions">
        <button className="btn primary" disabled={busy} onClick={handleSave}>
          {t.settings_save}
        </button>
        <button className="btn ghost" disabled={busy} onClick={onClose}>
          {t.settings_cancel}
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
  // Resolved locale for rendering: prefer the persisted choice from
  // state.json (returned by cmd_get_status), fall back to OS detection
  // before the first status arrives.
  const [locale, setLocale] = useState<Locale>(detectLocale());

  async function refresh() {
    try {
      const s = await invoke<StatusPayload>("cmd_get_status");
      setStatus(s);
      if (s.locale === "en" || s.locale === "zh") {
        setLocale(s.locale);
      }
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
          : { tray: event.payload, logged_in: true, cli_hits: [], base_url: "", locale: null }
      );
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

  async function saveLocale(next: Locale) {
    setLocale(next); // optimistic
    try {
      await invoke("cmd_set_locale", { locale: next });
    } catch (err) {
      console.error("cmd_set_locale failed", err);
    }
  }

  if (!status) {
    return (
      <main className="popup">
        <div className="loading">…</div>
      </main>
    );
  }

  const t = loadDict(locale);

  if (view === "settings") {
    return (
      <SettingsView
        baseUrl={status.base_url}
        locale={locale}
        t={t}
        busy={busy}
        onClose={() => setView("main")}
        onSave={saveServer}
        onLocaleChange={saveLocale}
      />
    );
  }

  return (
    <MainView
      status={status}
      t={t}
      busy={busy}
      onCommand={runCommand}
      onOpenSettings={() => setView("settings")}
      cliBannerDismissed={cliBannerDismissed}
      onDismissCliBanner={() => setCliBannerDismissed(true)}
    />
  );
}

export default App;
