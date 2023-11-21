import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {AiState} from "./ai.types";

const initialState: AiState = {
    searchContentText: '',
    usernames: '',
    groupFilter: '',
    workflowInstructions: '',
    searchResults: '',
    platformFilter: ''
}

const aiSlice = createSlice({
    name: 'ai',
    initialState,
    reducers: {
        setUsernames: (state, action: PayloadAction<string>) => {
            state.usernames = action.payload;
        },
        setGroupFilter: (state, action: PayloadAction<string>) => {
            state.groupFilter = action.payload;
        },
        setSearchContent: (state, action: PayloadAction<string>) => {
            state.searchContentText = action.payload;
        },
        setWorkflowInstructions: (state, action: PayloadAction<string>) => {
            state.workflowInstructions = action.payload;
        },
        setSearchResults: (state, action: PayloadAction<string>) => {
            state.searchResults = action.payload;
        },
        setPlatformFilter: (state, action: PayloadAction<string>) => {
            state.platformFilter = action.payload;
        }
    }
});

export const { setSearchContent,
    setGroupFilter,
    setUsernames,
    setWorkflowInstructions,
    setSearchResults,
    setPlatformFilter} = aiSlice.actions;
export default aiSlice.reducer;