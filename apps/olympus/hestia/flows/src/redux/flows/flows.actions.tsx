export interface FlowState {
    uploadContentContacts: any;
    uploadContentTasks: any,
    csvHeaders: string[];
    promptHeaders: string[];
}

export const initialState: FlowState = {
    uploadContentContacts: '',
    uploadContentTasks: '',
    csvHeaders: [],
    promptHeaders: []
}