const path = require('path');
const CopyPlugin = require("copy-webpack-plugin");
const HtmlWebpackPlugin = require("html-webpack-plugin");

module.exports = () => {
    return {
        target: ['web'],
        entry: './src/index.ts',
        module: {
            rules: [
                {
                    test: /\.ts?$/,
                    use: 'ts-loader',
                    exclude: /node_modules/,
                },
            ],
        },
        resolve: {
            extensions: ['.tsx', '.ts', '.js'],
        },
        output: {
            path: path.resolve(__dirname, 'dist'),
            filename: 'bundle.min.js',
            library: {
                type: 'umd'
            }
        },
        
        devServer: {
            static: path.join(__dirname, "dist"),
            compress: true,
            port: 3000,
        },
        plugins: [
            new CopyPlugin({
                // Use copy plugin to copy *.onnx to output folder.
                patterns: [{ from: 'src/*.onnx', to: '[name][ext]' }]
            }),
            new CopyPlugin({
                // Use copy plugin to copy *.jfif to output folder.
                patterns: [{ from: 'src/*.jfif', to: '[name][ext]' }]
            }),
            new CopyPlugin({
                // Use copy plugin to copy *.wasm to output folder.
                patterns: [{ from: 'node_modules/onnxruntime-web/dist/*.wasm', to: '[name][ext]' }]
            }),
            new HtmlWebpackPlugin({
                title: 'our project',
                template: 'src/index.html'
            })
        ],
        mode: 'production'
    }
};