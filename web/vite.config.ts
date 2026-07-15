import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import path from "path";
import tailwindcss from '@tailwindcss/vite'
import { tanstackRouter } from '@tanstack/router-plugin/vite'

export default defineConfig({
    plugins: [
        tanstackRouter({
            target: 'react',
            autoCodeSplitting: true,
            routesDirectory: './src/pages',  
            generatedRouteTree: './src/routeTree.gen.ts', // 生成的 routeTree 文件位置
        }),
        react(),
        tailwindcss(),
    ],
    server: {
        port: 3002,
    },
    resolve: {
        alias: {
            "@src": path.resolve(__dirname, "./src"),
        },
    },
});
