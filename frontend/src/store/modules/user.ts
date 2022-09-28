import { stat } from 'fs'
import { Module } from 'vuex'

export interface User {
  userName: string | null
  isAuthenticated: boolean
}

const defaultState: User = {
  userName: null,
  isAuthenticated: false
}

const userStore: Module<User, unknown> = {
  state: defaultState,
  mutations: {
    logIn(state, userName: string) {
      state.userName = userName
      state.isAuthenticated = true
    },
    logOut(state) {
      state = defaultState;
    },
  },
  actions: {
    logIn({ commit }, userName: string) {
      setTimeout(() => {
        commit('logIn', userName);
      }, 50);
    },
    logOut({ commit }) {
      setTimeout(() => {
        commit('logOut');
      }, 50);
    },
  }
}

export default userStore