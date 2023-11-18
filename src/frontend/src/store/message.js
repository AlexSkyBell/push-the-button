const state = {
    number: ""
}

const getters = {
    getNumber(state) {
        return state.number
    }
}

const mutations = {
    setNumber(state, number) {
        state.number = number
    }
}

const actions = {
    prepareMessage({ state, commit, getters }, response) {
        commit('setNumber', response)
    }
}

export default {
    state,
    getters,
    mutations,
    actions,
}