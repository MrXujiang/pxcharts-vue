import { resolve } from 'path'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import vueJsx from '@vitejs/plugin-vue-jsx'
import { TDesignResolver } from '@tdesign-vue-next/auto-import-resolver'

export default defineConfig({
  css: {
    preprocessorOptions: {
      // less 配置
      less: {
        modifyVars: {
          hack: `true; @import (reference) "${resolve(__dirname, 'src/assets/style/variables.less')}";`,
        },
        javascriptEnabled: true,
      },
    },
  },
  // 预构建依赖项，防止打包时重复打包，提升打包速度
  optimizeDeps: {
    include: ['vue', 'vue-router', 'pinia', '@tdesign-vue-next/auto-import-resolver'],
    exclude: ['vue-demi'],
  },
  // server: {
  //   port: 3190,
  //   proxy: {
  //     '/api': {
  //       target: 'https://test.admin.mvtable.com',
  //       changeOrigin: true,
  //       secure: false,
  //       rewrite: (path) => path,
  //     },
  //   },
  // },
  define: {
    'process.env': {
      BASE: '/mvtable/',
      BASE_API_URL: 'https://test.admin.mvtable.com/api/v1',
    },
  },
  plugins: [
    vue(),
    vueJsx(), // 支持 <script lang="tsx"> 和 JSX
    vueDevTools(),
    AutoImport({
      imports: ['vue', 'vue-router', 'pinia'], // 自动导入 Vue 和 Vue Router 的相关函数
      dts: true, // 生成对应的 .d.ts 文件
      dirs: [resolve(__dirname, './src/composables')], // 自定义 composables 目录
      resolvers: [
        TDesignResolver({
          library: 'vue-next',
        }),
      ],
    }),
    Components({
      resolvers: [
        TDesignResolver({
          library: 'vue-next',
        }),
      ],
    }),
  ],
  resolve: {
    alias: [
      {
        find: '@',
        replacement: resolve(__dirname, './src/'),
      },
    ],
  },
  base: '/mvtable/',
  build: {
    outDir: 'mvtable',
    minify: 'terser',
    rollupOptions: {
      output: {
        manualChunks: {
          vue: ['vue', 'vue-router', 'pinia'],
        },
      },
    },
    terserOptions: {
      compress: {
        //生产环境时移除console.log(), console.info, console.warn, console.error, 或者直接赋值true移除一切console.*的代码
        drop_console: false, // 禁用默认的移除所有console调用
        pure_funcs: ['console.log'], // 只移除console.log调用
        drop_debugger: true,
      },
    },
    chunkSizeWarningLimit: 2000,
  },
})
