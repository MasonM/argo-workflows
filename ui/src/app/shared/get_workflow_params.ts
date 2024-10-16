import {History} from 'history';

/**
 * Method to extract a key from query parameter. The query parameter will be of the following format
 * ?parameters[key]=value.
 * This method will extract the key from the query parameter.
 * @param inputString
 * @returns
 */
function extractKey(inputString: string): string | null {
    // Use regular expression to match the key within square brackets
    const match = inputString.match(/parameters\[(.*?)\]/);

    // If a match is found, return the captured key
    if (match) {
        return match[1];
    }

    // If no match is found, return null or an empty string
    return null; // Or return '';
}

export function getWorkflowParametersFromQuery(history: History): {[key: string]: string} {
    const queryParams = new URLSearchParams(history.location.search);
    const parameters: {[key: string]: string} = {};
    for (const [key, value] of queryParams.entries()) {
        const q = extractKey(key);
        if (q) {
            parameters[q] = value;
        }
    }

    return parameters;
}
