
export interface Assistant {
    id: string;
    object: string;
    created_at: number | null;
    name: string;
    description:  string | null;
    model: string;
    instructions?: string;
    tools?: any;
    file_ids?: any[];
    metadata?: any;
}

export interface DiscordFilters {
    categoryTopic?: string;
    categoryName?: string;
    category?: string;
}

export interface RetrievalItemInstruction {
    retrievalPlatform: string;
    retrievalPlatformGroups?: string;
    retrievalKeywords?: string;
    retrievalNegativeKeywords?: string;
    retrievalUsernames?: string;
    retrievalPrompt?: string;
    discordFilters?: DiscordFilters;
    webFilters?: WebFilters;
    instructions?: string;
}

export interface WebFilters {
    routingGroup?: string;
    lbStrategy?: string;
    maxRetries?: number;
    backoffCoefficient?: number;
}

export interface Retrieval {
    retrievalID?: number;
    retrievalName: string;
    retrievalGroup: string;
    retrievalItemInstruction: RetrievalItemInstruction;
}
