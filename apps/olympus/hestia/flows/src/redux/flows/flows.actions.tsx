export interface FlowState {
    uploadContentContacts: any;
    uploadContentTasks: any,
    csvHeaders: string[];
    promptHeaders: string[];
    results: any[];
}

export const initialState: FlowState = {
    uploadContentContacts: '',
    uploadContentTasks: '',
    csvHeaders: [],
    promptHeaders: [],
    results: []
}