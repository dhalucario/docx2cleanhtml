const path = require('path');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const WebpackNotifierPlugin = require('webpack-notifier');

module.exports = {
    mode: 'development',
    entry: './web/index.js',

    module: {
        rules: [
            {
                test: /\.js$/,
                exclude: /(node_modules|bower_components)/,
                use: {
                    loader: 'babel-loader',
                    options: {
                        presets: ['@babel/preset-env']
                    }
                }
            },
            {
                test: /\.scss$/,
                use: ExtractTextPlugin.extract({
                    disable: false,
                    allChunks: true,
                    fallback: 'style-loader',
                    use: ['css-loader', 'sass-loader']
                })
            },
            {
                test: /(htm|html)$/,
                use: [
                    'html-loader',
                    'htmllint-loader'
                ],
                exclude: /(node_modules)/,
            }
        ]
    },

    plugins: [
        new ExtractTextPlugin('bundle.css'),
        new WebpackNotifierPlugin(),
    ],

    output: {
        path: path.resolve('bin/public/'),
        filename: 'bundle.js'
    }
};
