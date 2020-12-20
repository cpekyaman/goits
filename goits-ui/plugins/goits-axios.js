export default function({$axios, redirect }) {
    $axios.onRequest(config => {
        console.log('sending request to ' + config.url)
    })
}