# Challenge

Crew is a recruiting CRM that helps companies manage their recruiting process. In order to get in touch with talents,
we offer to our customers to use enrichment platforms that, given a LinkedIn link, provide contacting data about the
talent.

These enrichment platforms have a cost per request. In order to bill back our customers, we have implemented a credit
system. Each workspace has a credit balance and each request to the enrichment platform costs a predefined number of
credits depending on the data requested (can be emails and phones).

Your goal will be to implement the credit management system.

## Requirements

Each workspace has a unique identifier and a name. By default, a workspace starts with 100 credits. One member from a
workspace can request data from the enrichment platform. Each enrichment, depending on the data that is requested, has a
cost: one credit for an email only request, two credits for an email and phone request. It is not possible to only get
the phone number.

We want you to implement the following features:

- A workspace can see their number of credits.
- A workspace can request data from the enrichment platform.
- A way, for Crew, to monitor the credits usage per workspace.

## Instructions

We do not expect you to build a way to create a workspace, users nor to authenticate them in any way.
For the sake of the test, you can just pass the workspaceID as part of a header.

We expect you to build a simple API that provides the two features above. You can use any language and framework you
want. For the storage part, we provide a postgres database in the `docker-compose.yml` file. You can use it to store the
workspaces and their credits. If you want to use any other database, feel free to do so.

For the enrichment platform, the service will be created when you run the `docker-compose up` command. It will be
available at `http://localhost:8080/api/v1/enrichment`. It has only one endpoint, `POST /api/v1/enrichment`.

```bash
# 200 OK
curl --location 'localhost:8080/api/v1/enrich' \
--header 'Content-Type: application/json' \
--data '{
    "linkedin_profile": "https://www.linkedin.com/in/steve-jobs"
}' \
-i

# 400 Bad Request
curl --location 'localhost:8080/api/v1/enrich' \
--header 'Content-Type: application/json' \
--data '{
    "linkedin_profile": ""
}' \
-i

# 404 Not Found
curl --location 'localhost:8080/api/v1/enrich' \
--header 'Content-Type: application/json' \
--data '{
    "linkedin_profile": "https://www.linkedin.com/in/steve-jobs-404"
}' \
-i
```

## Constraints

- A workspace cannot have a negative credit balance.
- A workspace cannot request data if it does not have enough credits.

## Deliverables

- An API that implements the features described above.
- A README file with instructions on how to run the project and how to interact with the API.
- A brief explanation of the architecture and design decisions.