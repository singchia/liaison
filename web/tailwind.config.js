/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/**/*.{js,jsx,ts,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: '#1890ff',
      },
    },
  },
  plugins: [],
  // 避免与 antd 样式冲突
  corePlugins: {
    preflight: false,
  },
};
