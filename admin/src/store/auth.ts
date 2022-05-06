import { defineStore } from 'pinia';

interface AuthStoreState {
  isLogged: boolean;
}

export const useAuthStore = defineStore('authStore', {
  state: () => ({ isLogged: false }),

  getters: {
    getAuthState: (state) => {
      return state.isLogged;
    },
  },

  actions: {
    updateAuthState(authState: boolean) {
      this.isLogged = authState;
    },
  },
});
