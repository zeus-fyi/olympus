export interface FlowState {
    uploadContent: any;
    csvHeaders: string[];
}

export const initialState: FlowState = {
    uploadContent: '',
    csvHeaders: []
}