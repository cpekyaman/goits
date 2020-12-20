export const state = () => ({
    item: {}
})
  
export const mutations = {
    item(state, data) {
        state.item = data
    }
}

export const getters = {
    item(state) {
        return state.item
    }
}