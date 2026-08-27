import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

import { tokenStorage } from "@/shared/api/tokenStorage";
import type { Session, User } from "@/shared/types";

const USER_KEY = "platforma.user";

function restoreUser(): User | null {
  try {
    const raw = localStorage.getItem(USER_KEY);
    return raw && tokenStorage.access() ? (JSON.parse(raw) as User) : null;
  } catch {
    return null;
  }
}

type AuthState = {
  user: User | null;
  initialized: boolean;
};

const initialState: AuthState = {
  user: restoreUser(),
  initialized: false,
};

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    sessionStarted(state, action: PayloadAction<Session>) {
      const { accessToken, refreshToken, user } = action.payload;
      tokenStorage.save(accessToken, refreshToken);
      localStorage.setItem(USER_KEY, JSON.stringify(user));
      state.user = user;
      state.initialized = true;
    },
    userRefreshed(state, action: PayloadAction<User>) {
      localStorage.setItem(USER_KEY, JSON.stringify(action.payload));
      state.user = action.payload;
      state.initialized = true;
    },
    sessionEnded(state) {
      tokenStorage.clear();
      localStorage.removeItem(USER_KEY);
      state.user = null;
      state.initialized = true;
    },
    authInitialized(state) {
      state.initialized = true;
    },
  },
});

export const { sessionStarted, userRefreshed, sessionEnded, authInitialized } = authSlice.actions;
export default authSlice.reducer;
