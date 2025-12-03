/**
 * 在生产环境 代理是无法生效的，所以这里没有生产环境的配置
 * -------------------------------
 * The agent cannot take effect in the production environment
 * so there is no configuration of the production environment
 * For details, please see
 * https://pro.ant.design/docs/deploy
 */
export default {
  dev: {
    '/api': {
      target: 'http://82.156.45.8',
      secure: false,
      changeOrigin: true,
      pathRewrite: { '^/api': '/api' },
    },
  },
  test: {
    '/api': {
      target: 'https://192.168.119.200',
      secure: false,
      changeOrigin: true,
      pathRewrite: { '^/api': '/api' },
    },
  },
};
