export default {
    methods: {
        getLabel(group, suffix) {
            let customRoot = null
            if(this.$parent && this.$parent.rootPath) {
                customRoot = this.$parent.rootPath
            } else if(this.rootPath) {
                customRoot = this.rootPath
            }

            let labelResource = ''
            if(customRoot) {
                labelResource = this.$i18n.t(customRoot + '.' + group + '.' + suffix)
            }
            if(labelResource === '') {
                labelResource = this.$i18n.t(group + '.' + suffix)
            }
            return labelResource
        }
    }
}