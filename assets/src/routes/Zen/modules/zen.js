import axios from 'axios'

// ------------------------------------
// Constants
// ------------------------------------
export const REQUEST_ZEN='REQUEST_ZEN';
export const RECEIVE_ZEN='RECEIVE_ZEN';
export const SAVE_CURRENT_ZEN='SAVE_CURRENT_ZEN';

// ------------------------------------
// Actions
// ------------------------------------
export function requestZen () {
    return {
        type: REQUEST_ZEN
    }
}

let availableId = 0;
export function receiveZen(value) {
    return {
        type: RECEIVE_ZEN,
        payload: {
            value,
            id: availableId++
        }
    }
}

export function saveCurrentZen() {
    return {
        type: SAVE_CURRENT_ZEN
    }
}

export const fetchZen = () => {
    return (dispatch) => {
        dispatch(requestZen());
        return axios.get('https://api.github.com/zen')
            .then(response => dispatch(receiveZen(response.data)))
    }
};

export const actions = [
    requestZen,
    receiveZen,
    fetchZen,
    saveCurrentZen
];

// ------------------------------------
// Action Handlers
// ------------------------------------
const ACTION_HANDLERS = {
    [REQUEST_ZEN]: (state) => ({ ...state, fetching: true }),
    [RECEIVE_ZEN]: (state, action) => ({
        ...state,
        fetching: false,
        zens: state.zens.concat(action.payload),
        current:action.payload.id
    }),
    [SAVE_CURRENT_ZEN]: (state) =>
        state.current != null ?
            ({...state, saved: state.saved.concat(state.current)}) :
            state
};

// ------------------------------------
// Reducer
// ------------------------------------
const initialState = {fetching: false, current:null, zens: [], saved: []};
export default function zenReducer (state = initialState, action) {
    const handler = ACTION_HANDLERS[action.type];
    return handler ? handler(state, action) : state
}

