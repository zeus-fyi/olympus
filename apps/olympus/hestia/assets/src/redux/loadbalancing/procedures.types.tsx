type Queue = any[];

interface IrisRoutingProcedure {
    name: string;
    orderedSteps?: Queue;
}

interface BroadcastInstructions {
    routingPath: string;
    restType: string;
    payload?: any;
    maxDuration?: number; // time.Duration equivalent
    maxTries?: number;
    routingTable: string;
    fanInRules?: FanInRules;
}

interface FanInRules {
    rule: BroadcastRules;
}

type BroadcastRules = "returnOnFirstSuccess" | "returnAllSuccessful";

const FanInRuleFirstValidResponse: BroadcastRules = "returnOnFirstSuccess";
const FanInRuleReturnAllResponses: BroadcastRules = "returnAllSuccessful";

// ReturnFirstResultOnSuccess function
function returnFirstResultOnSuccess(b: BroadcastRules): BroadcastRules {
    return FanInRuleFirstValidResponse;
}

// ReturnResultsOnSuccess function
function returnResultsOnSuccess(b: BroadcastRules): BroadcastRules {
    return FanInRuleReturnAllResponses;
}

// Assuming iris_operators.IrisRoutingResponseETL and iris_operators.Aggregation are already defined
interface IrisRoutingProcedureStep {
    broadcastInstructions?: BroadcastInstructions;
    transformSlice?: iris_operators.IrisRoutingResponseETL[]; // Replace with actual type
    aggregateMap?: { [key: string]: iris_operators.Aggregation }; // Replace with actual type
}

// Assume iris_operators is imported or defined
namespace iris_operators {
    export interface IrisRoutingResponseETL {
        // define your properties and methods here
    }

    export interface Aggregation {
        // define your properties and methods here
    }
}