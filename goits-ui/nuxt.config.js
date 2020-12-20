export default {
    plugins: [
        '~/plugins/goits-components.js',
        '~/plugins/goits-axios.js'
    ],

    modules: [
        ['nuxt-i18n', { 
            locales: [ 
                { code: 'en', iso: 'en-US', file: 'en.json' },
                { code: 'tr', iso: 'tr-TR', file: 'tr.json' },
            ],
            defaultLocale: 'en',
            langDir: './locales/',
            lazy: true,
         },
        ],

        ['@nuxtjs/axios', {
            prefix: '/api/',
            proxy: true,
            headers: {
                common: {
                    'Accept': 'application/json, text/plain, */*'
                }
            }
        }],
    ],

    buildModules : [
        ['@nuxtjs/vuetify', { optionsPath: './vuetify.options.js' }],
    ],

    proxy: {
        '/api': { 
            target: 'http://localhost:8080', 
            pathRewrite: {'^/api': ''} 
        }
    }
}