<template>
    <g-entity-list :columns="columns" :itemKey="itemKey" :editable="true" :removable="true" />
</template>

<script>
export default {
    data() {
        return {
            columns: ['id','name'],
            itemKey: "id"
        }
    }
}
</script>