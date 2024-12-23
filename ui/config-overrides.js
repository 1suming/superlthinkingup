/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

const {
  addWebpackModuleRule,
  addWebpackAlias,
  setWebpackOptimizationSplitChunks,
} = require("customize-cra");

const path = require("path");
const i18nPath = path.resolve(__dirname, "../i18n");

// const UnoCSS = require('@unocss/webpack').default


module.exports = {
  webpack: function(config, env) {
    console.log("process.env.NODE_ENV:",process.env.NODE_ENV)
    console.log("@process.env.REACT_APP_API_URL:", process.env.REACT_APP_API_URL)
    addWebpackAlias({
      "@": path.resolve(__dirname, "src"),
      "@i18n": i18nPath
    })(config);

    addWebpackModuleRule({
      test: /\.ya?ml$/,
      use: "yaml-loader"
    })(config);

    setWebpackOptimizationSplitChunks({
      maxInitialRequests: 20,
      minSize: 20 * 1024,
      minChunks: 2,
      cacheGroups: {
        automaticNamePrefix: 'chunk',
        mix1: {
          test: (module, chunks) => {
            return (
              module.resource &&
              (module.resource.includes('components') ||
                /\/node_modules\/react-bootstrap\//.test(module.resource))
            );
          },
          name: 'chunk-mix1',
          filename: 'static/js/[name].[contenthash:8].chunk.js',
          priority: 14,
          reuseExistingChunk: true,
          minChunks: process.env.NODE_ENV === 'production' ? 1 : 2,
          chunks: 'initial',
        },
        mix2: {
          name: 'chunk-mix2',
          test: /[\\/]node_modules[\\/](i18next|lodash|marked|next-share)[\\/]/,
          filename: 'static/js/[name].[contenthash:8].chunk.js',
          priority: 13,
          reuseExistingChunk: true,
          minChunks: 1,
          chunks: 'initial',
        },
        mix3: {
          name: 'chunk-mix3',
          test: /[\\/]node_modules[\\/](@remix-run|@restart|axios|diff)[\\/]/,
          filename: 'static/js/[name].[contenthash:8].chunk.js',
          priority: 12,
          reuseExistingChunk: true,
          minChunks: 1,
          chunks: 'initial',
        },
        codemirror: {
          name: 'codemirror',
          test: /[\\/]node_modules[\\/](\@codemirror)[\\/]/,
          priority: 10,
          reuseExistingChunk: true,
          minChunks: process.env.NODE_ENV === 'production' ? 1 : 2,
          chunks: 'initial',
          enforce: true,
        //   filename: 'static/js/[name].[contenthash:8].chunk.js',
        },
        lezer: {
          name: 'lezer',
          test: /[\\/]node_modules[\\/](\@lezer)[\\/]/,
          priority: 9,
          reuseExistingChunk: true,
          minChunks: process.env.NODE_ENV === 'production' ? 1 : 2,
          chunks: 'initial',
          enforce: true,
          filename: 'static/js/[name].[contenthash:8].chunk.js', //@my add
        },
        reactDom: {
          name: 'react-dom',
          test: /[\\/]node_modules[\\/](react-dom)[\\/]/,
          filename: 'static/js/[name].[contenthash:8].chunk.js',
          priority: 8,
          reuseExistingChunk: true,
          chunks: 'all',
          enforce: true,
        },
        nodesInitial: {
          name: 'chunk-nodesInitial',
          filename: 'static/js/[name].[contenthash:8].chunk.js',
          test: /[\\/]node_modules[\\/]/,
          priority: 1,
          minChunks: 1,
          chunks: 'initial',
          reuseExistingChunk: true,
        },

        tinymce: {
            name: 'tinymce',
            test: /[\\/]node_modules[\\/](\@tinymce)[\\/]/,
            priority: 30,
            reuseExistingChunk: true,
            minChunks: process.env.NODE_ENV === 'production' ? 1 : 2,
            chunks: 'initial',
            enforce: true,
          },

      },
    })(config);

    // add i18n dir to ModuleScopePlugin allowedPaths
    const moduleScopePlugin = config.resolve.plugins.find(_ => _.constructor.name === "ModuleScopePlugin");
    if (moduleScopePlugin) {
      moduleScopePlugin.allowedPaths.push(i18nPath);
    }

    const BundleAnalyzerPlugin = require('webpack-bundle-analyzer').BundleAnalyzerPlugin;
    if (process.env.ANALYZE) {
        config.plugins.push(new BundleAnalyzerPlugin({

            // analyzerMode: 'static',
           

              reportFilename: 'report.html',
              openAnalyzer: true, // 自动打开浏览器
            }
            
        ));
  
   } //else {
//     //     config.plugins = config.plugins.filter(plugin => plugin.constructor.name !== 'BundleAnalyzerPlugin');
//     // }
/*
打包不生成 source-map 文件
在 配置文件 config-overrides.js 添加 process.env.GENERATE_SOURCEMAP = "false";
或者在项目更目录下创建文件 .env, 写入: GENERATE_SOURCEMAP=false 即可.
*/

    // this.addWebpackPlugin("UnoCSS", new UnoCSS({ configFile: './unocss.config.js' }))
    // config.optimization.realContentHash = true;
    import('@unocss/webpack').then(({ default: UnoCSS }) => ({
        plugins: [UnoCSS()],
        optimization: {
          realContentHash: true
        }
      }));
 
    return config;
  },
  devServer: function(configFunction) {
    console.log("@@process.env.REACT_APP_API_URL",process.env.REACT_APP_API_URL)
    return function(proxy, allowedHost) {
      const config = configFunction(proxy, allowedHost);
      config.proxy = [
        {
          context: ['/answer', '/installation'], //,'/article'
          target: process.env.REACT_APP_API_URL,
          changeOrigin: true,
          secure: false,
        },
        {
          context: ['/custom.css'],
          target: process.env.REACT_APP_API_URL,
        }
      ];
      return config;
    };
  }
};
