// Tiny i18n for the desktop popup. Two locales for now: zh and en.
//
// The shape is deliberately a flat key-value dictionary, not a nested
// namespace tree — at this scale (one popup, ~30 strings) the lookup
// cost of an extra branch buys nothing. New keys go in the `Dict`
// type so missing translations fail at compile time rather than
// silently fall back to the key string.

export type Locale = "en" | "zh";

export interface Dict {
  // Status pill
  label_logged_out: string;
  label_paused: string;
  label_connecting: string;
  label_online: string;
  label_error: string;

  // Action buttons
  btn_login: string;
  btn_pause: string;
  btn_resume: string;
  btn_reconnect: string;
  btn_dashboard: string;
  btn_logout: string;

  // CLI banner (first-run on machines that already have liaison-edge)
  cli_banner_title: string;
  cli_banner_body: string;
  cli_banner_more: (n: number) => string;
  cli_banner_dismiss: string;

  // Settings page
  settings_back_aria: string;
  settings_open_aria: string;
  settings_title: string;
  settings_url_label: string;
  settings_url_placeholder: string;
  settings_url_hint: string;
  settings_save: string;
  settings_cancel: string;
  settings_url_empty: string;

  // Locale picker
  settings_locale_label: string;
  locale_en: string;
  locale_zh: string;
}

const en: Dict = {
  label_logged_out: "Not signed in",
  label_paused: "Paused",
  label_connecting: "Connecting",
  label_online: "Connected",
  label_error: "Error",

  btn_login: "Sign in to Liaison",
  btn_pause: "Pause",
  btn_resume: "Resume",
  btn_reconnect: "Reconnect",
  btn_dashboard: "Open Dashboard",
  btn_logout: "Sign out",

  cli_banner_title: "⚠ An existing liaison-edge CLI install was detected",
  cli_banner_body:
    "Continuing with Liaison Desktop will create a separate connector that runs alongside the CLI. Consider uninstalling the CLI version first:",
  cli_banner_more: (n) => `+ ${n} more…`,
  cli_banner_dismiss: "Got it, continue",

  settings_back_aria: "Back",
  settings_open_aria: "Settings",
  settings_title: "Server settings",
  settings_url_label: "Liaison server URL",
  settings_url_placeholder: "https://liaison.example.com",
  settings_url_hint:
    "Public users keep the default https://liaison.cloud. For a private deployment, enter your URL starting with http:// or https://.",
  settings_save: "Save and re-sign in",
  settings_cancel: "Cancel",
  settings_url_empty: "Address cannot be empty",

  settings_locale_label: "Language",
  locale_en: "English",
  locale_zh: "中文",
};

const zh: Dict = {
  label_logged_out: "未登录",
  label_paused: "已暂停",
  label_connecting: "连接中",
  label_online: "已连接",
  label_error: "错误",

  btn_login: "登录 Liaison",
  btn_pause: "暂停连接",
  btn_resume: "恢复连接",
  btn_reconnect: "重连",
  btn_dashboard: "打开 Dashboard",
  btn_logout: "退出登录",

  cli_banner_title: "⚠ 检测到本机已有 liaison-edge CLI 安装",
  cli_banner_body:
    "继续使用 Liaison Desktop 会创建一个独立的连接器，与 CLI 并行运行。建议先卸载 CLI 版本：",
  cli_banner_more: (n) => `+ ${n} 项…`,
  cli_banner_dismiss: "知道了，继续",

  settings_back_aria: "返回",
  settings_open_aria: "设置",
  settings_title: "服务器配置",
  settings_url_label: "Liaison 服务器地址",
  settings_url_placeholder: "https://liaison.example.com",
  settings_url_hint:
    "公网用户保持默认 https://liaison.cloud。私有化部署请填写你的部署地址，以 http:// 或 https:// 开头。",
  settings_save: "保存并重新登录",
  settings_cancel: "取消",
  settings_url_empty: "地址不能为空",

  settings_locale_label: "语言",
  locale_en: "English",
  locale_zh: "中文",
};

const dicts: Record<Locale, Dict> = { en, zh };

/// Pick a locale based on whatever the OS / browser exposes via
/// navigator.language. Anything starting with "zh" is treated as
/// Simplified Chinese; everything else falls back to English.
export function detectLocale(): Locale {
  if (
    typeof navigator !== "undefined" &&
    typeof navigator.language === "string" &&
    navigator.language.toLowerCase().startsWith("zh")
  ) {
    return "zh";
  }
  return "en";
}

export function dict(locale: Locale): Dict {
  return dicts[locale] ?? dicts.en;
}
