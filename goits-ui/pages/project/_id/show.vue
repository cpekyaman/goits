<template>
    <g-entity-detail :itemId="$route.params.id">
        <g-display-text :dataModel="item" field="id" />
        <g-display-text :dataModel="item" field="name" />
    </g-entity-detail>
</template>

<script>
import { mapGetters } from "vuex";

export default {
    data() {
        return {
            id: 0
        }
    },

    computed: {
         ...mapGetters("project", [
            "item",
        ])
    },

    validate({params}) {
        return /^\d+$/.test(params.id)
    }
}
</script>