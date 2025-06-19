import {ElasticResponse} from "./models/elastic";

const baseUrl = 'https://localhost:8443';
const indexName = "recipe-problems";

export async function QueryReports(startAt: number, pageSize: number):Promise<ElasticResponse> {
    const url = `${baseUrl}/${indexName}/_search`;
    const query = {
        match_all: {}
    };

    const req = {
        from: startAt,
        size: pageSize,
        query,
        sort: [
            {
                annotation_count: {
                    order: "desc"
                }
            }
        ]
    };

    const response = await fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(req)
    });
    if(response.status!=200) {
        const responseText = await response.text();
        throw new Error('server responded' + responseText);
    }

    return ElasticResponse.parse(await response.json());
}