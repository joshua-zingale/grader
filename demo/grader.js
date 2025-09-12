/**
 * A response from an activity submission.
 * @typedef {Object} SubmissionFeedback
 * @property {number} grade - the evaluation for the submission, 0 <= grade <= 1
 * @property {string} hint - the hint given for the submitted answer; an empty string if there is no hint.
 */

/**
 * Submits an activity to the grading server and returns
 * @param {string} graderUrl - The URL to the activity grading web server.
 * @param {string} identifier - The unique activity ID for which an answer is to be submitted.
 * @param {string} answer - The answer to be submitted for grading.
 * @param {string} session - The. nique session ID for a user
 * @returns {Promise<SubmissionFeedback>} A promise that resolves to the graded feedback.
 * @throws {Error} Throws an error if the network request fails or the server returns a non-200 status.
 */
async function gradeSubmission(graderUrl, identifier, answer, session) {
  return fetch(graderUrl, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      identifier,
      answer,
      session,
    }),
  }).then(async (response) => {
    if (!response.ok) {
      throw Error(await response.text());
    }
    return response.json();
  }).catch((error) => {
    throw error;
  });
}