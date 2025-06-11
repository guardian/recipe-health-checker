import {z} from "zod";
import {CheckerOutput} from "./recipe-health-checker.ts";

export const SingleHitResponse = z.object({
    _index: z.string(),
    _id: z.string(),
    _score: z.number().nullable(),
    _source: CheckerOutput,
    sort: z.array(z.union([z.string(), z.number()])).optional().nullable()
});

const ValueResponse = z.object({
    value: z.number(),
    relation: z.string()    //should be enum, meh
});

const HitsResponse =  z.object({
    total: ValueResponse,
    max_score: z.number().optional().nullable(),
    hits: z.array(SingleHitResponse)
});

export const ElasticResponse = z.object({
    took: z.number(),
    timed_out: z.boolean(),
    hits: HitsResponse
});
export type ElasticResponse = z.infer<typeof ElasticResponse>;
export type SingleHitResponse = z.infer<typeof SingleHitResponse>;
