import BaseComponent from '@/components/mixin/BaseComponent.js'

export default {
    mixins:[BaseComponent],

    props: {
        field: {type: String, required: true},
        dataModel: {type: Object, required: true},
        label: {type: String},
    },

    computed: {
        inputLabel() {
            if(this.label) {
                return label
            }
            return this.getLabel('form.label', this.field)
        }
    }
}