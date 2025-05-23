import http from "k6/http";
import {check} from "k6";
import encoding from "k6/encoding";

export const options = {
    // iterations: 1,
    stages: [
        { duration: '1m', target: 3 },
        { duration: '3m', target: 12 },
        { duration: '7m', target: 24 },
        { duration: '5m', target: 10 },
        { duration: '3m', target: 0 }
    ],
};

export default function() {
    const workflowURL = "http://localhost:5000/api/workflow";
    const params = {
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Basic "+encoding.b64encode("admin:qwerty")
        }
    };
    const payload = JSON.stringify({
        name: "merchant_callback",
        version: 1,
        priority: 1,
        input: {
            merchantId: "dc1754de-0c00-4abc-afd8-b10ffe061024",
            event: "callback_test",
            payload: {
                "message": "test"
            },
            additionalInfo: {},
        },
    });

    const workflowExecResp = http.post(workflowURL, payload, params);
    check(workflowExecResp, { 'execute status was 200': (r) => r.status == 200 });


    const getWorkflowURL = "http://localhost:5000/api/workflow/"+workflowExecResp.body+"?includeTasks=true";
    const workflowDataResp = http.get(getWorkflowURL);
    check(workflowDataResp, { 'workflow result status was 200': (r) => r.status == 200})
}