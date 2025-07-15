import type {CheckerOutput} from "./models/recipe-health-checker.ts";

export interface Annotation {
    section: string;
    detail: string;
}

export interface AnnotationSummary {
    content: Annotation[];
    count: number;
    breakdown: Record<string, number>;
}

function get_section(report:CheckerOutput, idx:number):string {
    if(report.annotation_section && report.annotation_section[idx]) {
        return report.annotation_section[idx].substring(1)
    } else {
        return "unknown"
    }
}
function get_summary(report:CheckerOutput, idx:number):string {
    if(report.annotation_summary && report.annotation_summary[idx]) {
        return report.annotation_summary[idx]
    } else {
        return "unknown"
    }
}

export function GetAnnotationSummary(report:CheckerOutput):AnnotationSummary {
    const annotations:Annotation[] = [];

    for(let i=0;i<report.annotation_count;i++) {
        annotations.push({
            section: get_section(report, i),
            detail: get_summary(report, i),
        })
    }

    const breakdownMap:Map<string, number> = new Map();
    for(let i=0;i<report.annotation_count;i++) {
        const s = get_section(report, i);
        breakdownMap.set(s, (breakdownMap.get(s) ?? 0)+1)
    }
    return {
        content: annotations,
        count:report.annotation_count,
        breakdown: Object.fromEntries(breakdownMap)
    }
}