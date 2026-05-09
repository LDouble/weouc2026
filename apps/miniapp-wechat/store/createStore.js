function cloneState(value) {
  if (value === undefined) return undefined;
  return JSON.parse(JSON.stringify(value));
}

export default function createStore(initialState = {}) {
  let state = cloneState(initialState);
  let listeners = [];

  function getState() {
    return cloneState(state);
  }

  function setState(updater, replace = false) {
    const previousState = getState();
    const nextState = typeof updater === 'function' ? updater(previousState) : updater;

    if (!nextState || typeof nextState !== 'object') {
      return getState();
    }

    state = replace ? cloneState(nextState) : Object.assign({}, state, cloneState(nextState));

    const currentState = getState();
    listeners.slice().forEach((listener) => {
      listener(currentState, previousState);
    });

    return currentState;
  }

  function subscribe(listener) {
    if (typeof listener !== 'function') return () => {};

    listeners.push(listener);

    return () => {
      listeners = listeners.filter((item) => item !== listener);
    };
  }

  function reset() {
    return setState(initialState, true);
  }

  return {
    getState,
    setState,
    subscribe,
    reset,
  };
}
