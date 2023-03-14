const environment = require('./environment');
const paths = require('./paths');
const webpack = require('webpack');
const path = require('path');
const { WebpackManifestPlugin } = require('webpack-manifest-plugin');

const plugins = {
    Clean: require('clean-webpack-plugin').CleanWebpackPlugin,
    Define: require('webpack').DefinePlugin,
    Html: require('html-webpack-plugin'),
    MiniCSSExtract: require('mini-css-extract-plugin'),
    Size: require('size-plugin'),
    Copy: require('copy-webpack-plugin'),
    PrerenderSPAPlugin: require('prerender-spa-plugin'),
};

module.exports = {
    common: [],
    start: [
        new plugins.Define({
            'process.env': environment.development,
        }),
        new webpack.ProvidePlugin({
            Buffer: ['buffer', 'Buffer'],
        }),
        new plugins.Html({
            inject: true,
            template: paths.APP_HTML,
        }),
    ],
    build: [
        new plugins.Clean(),
        new plugins.Copy({
            patterns: [
                {
                    from: 'public/*.(svg|png|jpg|jpeg)',
                    to: '[name][ext]',
                },
                {
                    from: 'public/images/*.(svg|png|jpg|jpeg)',
                    to: 'images/[name][ext]',
                },
            ],
        }),
        new plugins.MiniCSSExtract({
            filename: 'static/css/[name].[chunkhash:8].css',
            chunkFilename: 'static/css/[id].[chunkhash:8].css',
        }),
        new plugins.Html({
            inject: 'body',
            template: paths.APP_HTML,
            minify: false,
        }),
        new plugins.Size(),
        new WebpackManifestPlugin(),
    ],
};
