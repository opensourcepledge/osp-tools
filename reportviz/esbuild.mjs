// © 2024 Vlad-Stefan Harbuz <vlad@vlad.website>
// SPDX-License-Identifier: Apache-2.0

import * as esbuild from 'esbuild';
import ImportGlobPlugin from 'esbuild-plugin-import-glob';
const ImportGlob = ImportGlobPlugin.default

const ctx = await esbuild.context({
    entryPoints: ['src/index.js'],
    outfile: 'build/index.js',
    bundle: true,
    plugins: [
        ImportGlob(),
    ],
});

await ctx.watch();

const { host, port } = await ctx.serve({
    servedir: '.',
});

console.log(`Listening on ${host} port ${port}`);
