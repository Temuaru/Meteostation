const path = require('path')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin'); 

module.exports = {
  mode: 'development', 
  devtool: 'inline-source-map',
  entry: './dashboard/src/scripts/index.js',
  module: {
    rules: [
      { test: /\.svg$/, use: 'svg-inline-loader' },
      { test: /\.css$/, use: [ MiniCssExtractPlugin.loader, 'css-loader' ] },
      { 
        test: /\.(js)$/, 
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
          options: { presets: ['@babel/preset-env'] }
        }
      }
    ]
  },
  output: {
    path: path.resolve(__dirname, 'dashboard', 'dist'),
    filename: 'index_bundle.js',
    publicPath: '/dist',
    clean: true
  }, 
  plugins: [
    new HtmlWebpackPlugin({
      template: './dashboard/src/index.html',
    }),
    new MiniCssExtractPlugin({
        filename: 'style.css',
    }),
  ]
}
