const { defineConfig } = require('@vue/cli-service')
const webpack = require("webpack");
module.exports = defineConfig({
  transpileDependencies: true,
  configureWebpack: {
    plugins: [
      // janus.js does not use 'import' to access to the functionality of webrtc-adapter,
      // instead it expects a global object called 'adapter' for that.
      // Let's make that object available.
      new webpack.ProvidePlugin({ adapter: ['webrtc-adapter', 'default'] })
    ]
  }
})
