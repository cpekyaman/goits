export default {
    data() {
        return {
            rootPath: ''
        }
    }, 

    created() {
        let path = this.$route.path.substring(1)
        const subpathPos = path.indexOf('/')
        if(subpathPos > 0) {
            this.rootPath = path.substring(0, subpathPos)
        } else {
            this.rootPath = path.substring(0)
        }
        console.log(this.rootPath)
    }
}