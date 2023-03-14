// eslint-disable-next-line import/no-unresolved
// eslint-disable-next-line import/extensions,import/no-unresolved
const path = require('path');
// eslint-disable-next-line import/no-unresolved
const sveltePreprocess = require('svelte-preprocess');
// eslint-disable-next-line import/extensions,import/no-unresolved
const {paths, plugins} = require('./utils');

module.exports = {
  entry: {
    app: paths.APP_ENTRY_POINT,
  },
  devServer: {
    historyApiFallback: true,
  },
  module: {
    rules: [
      {
        test: /\.svelte$/,
        use: {
          loader: 'svelte-loader',
          options: {
            emitCss: true,
            hotReload: true,
            preprocess: sveltePreprocess({
              postcss: true,
            }),
          },
        },
      },
      {
        test: /\.m?js$/,
        use: {
          loader: 'babel-loader',
        },
      },
      {
        test: /\.(png|jpe?g|svg)$/i,
        use: {
          loader: 'file-loader',
        },
      },
    ],
  },
  resolve: {
    alias: {
      svelte: path.resolve('node_modules', 'svelte'),
      ws: path.resolve(__dirname, '../src/ws'),
      js: path.resolve(__dirname, '../src/js'),
      Components: path.resolve(__dirname, '../src/components/'),
      Layouts: path.resolve(__dirname, '../src/layouts/'),
    },
    fallback: {
      buffer: require.resolve('buffer'),
    },
    extensions: ['.mjs', '.js', '.svelte'],
    mainFields: ['svelte', 'browser', 'module', 'main'],
  },
  plugins: plugins.common,
};
