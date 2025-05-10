import { test, expect } from "@playwright/test";
import { faker } from "@faker-js/faker";

let email, authToken, feed_id, url, url2, feed_id2, email2, authToken2;
test.beforeEach("Credentials - User", async ({ request }) => {
  email = faker.internet.email();
  const password = faker.internet.password();
  const name = faker.person.firstName();
  email2 = faker.internet.email();
  // Create User 1
  const response = await request.post("/v1/user", {
    data: {
      email: email,
      password: password,
      name: name,
    },
  });
  // Create user 2
  const response2 = await request.post("/v1/user", {
    data: {
      email: email2,
      password: password,
      name: name,
    },
  });
  // Login User 1
  const loginResponse = await request.post("/v1/login", {
    form: {
      username: email,
      password: password,
    },
  });
  // Login User 2
  const loginResponse2 = await request.post("/v1/login", {
    form: {
      username: email2,
      password: password,
    },
  });
  const json2 = await loginResponse2.json();
  authToken2 = json2.token;
  const json = await loginResponse.json();
  authToken = json.token;

  // Validate status code
  expect(loginResponse.status()).toBe(200);
  expect(response.status()).toBe(201);
  expect(response2.status()).toBe(201);
  expect(loginResponse2.status()).toBe(200);
});

test.beforeEach("Credentials - Feed", async ({ request }) => {
  url = faker.internet.url();
  url2 = faker.internet.url();
  const response = await request.post("/v2/feeds", {
    headers: {
      Authorization: `Bearer ${authToken}`,
    },
    data: {
      name: faker.lorem.word(),
      url: url,
    },
  });
  const json = await response.json();
  feed_id = json.id;
  expect(response.status()).toBe(201);

  const response2 = await request.post("/v2/feeds", {
    headers: {
      Authorization: `Bearer ${authToken}`,
    },
    data: {
      name: faker.lorem.word(),
      url: url2,
    },
  });
  const json2 = await response2.json();
  feed_id2 = json2.id;
  expect(response2.status()).toBe(201);
});

test.afterEach("Remove Credentials", async ({ request }) => {
  // Remove feed 1
  await request.delete(`/v2/feeds/${feed_id}`, {
    headers: {
      Authorization: `Bearer ${authToken}`,
    },
  });
  // Remove feed 2
  await request.delete(`/v2/feeds/${feed_id2}`, {
    headers: {
      Authorization: `Bearer ${authToken}`,
    },
  });
  // Remove user 2
  const res = await request.delete(`/v1/user`, {
    headers: {
      Authorization: `Bearer ${authToken2}`,
    },
  });
  // Remove user 1
  const res2 = await request.delete("/v1/user", {
    headers: {
      Authorization: `Bearer ${authToken}`,
    },
  });
  expect(res.status()).toBe(204);
  expect(res2.status()).toBe(204);
});

test.describe("Create Feed", () => {
  test("Create Feed", async ({ request }) => {
    const name = faker.lorem.word();
    const url2 = faker.internet.url();
    const response = await request.post("/v2/feeds", {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
      data: {
        name: name,
        url: url2,
      },
    });

    // Validate status code
    expect(response.status()).toBe(201);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("name", name);
    expect(json).toHaveProperty("url", url2);
    expect(json).toHaveProperty("id");
    await request.delete(`/v2/feeds/${json.id}`, {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    });
  });

  test("Create Feed - Not Authorized", async ({ request }) => {
    const response = await request.post("/v2/feeds", {
      data: {
        name: faker.lorem.word(),
        url: faker.internet.url(),
      },
    });
    // Validate status code
    expect(response.status()).toBe(401);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Unauthorized");
  });

  test("Create Feed - URL Exsits", async ({ request }) => {
    const response = await request.post("/v2/feeds", {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
      data: {
        name: faker.lorem.word(),
        url: url,
      },
    });
    // Validate status code
    expect(response.status()).toBe(409);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Feed exist");
  });

  test("Create Feed - Invalid URL", async ({ request }) => {
    const response = await request.post("/v2/feeds", {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
      data: {
        name: faker.lorem.word(),
        url: "Invalid URL",
      },
    });
    // Validate status code
    expect(response.status()).toBe(400);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Invalid URL");
  });
});

test.describe("Update Feed", () => {
  test("Update Feed", async ({ request }) => {
    const name = faker.lorem.word();
    const url2 = faker.internet.url();
    const response = await request.put(`/v2/feeds/${feed_id}`, {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
      data: {
        name: name,
        url: url2,
      },
    });

    // Validate status code
    expect(response.status()).toBe(200);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("name", name);
    expect(json).toHaveProperty("url", url2);
    expect(json).toHaveProperty("id", feed_id);
  });

  test("Update Feed - Not Authorized", async ({ request }) => {
    const response = await request.put(`/v2/feeds/${feed_id}`, {
      data: {
        name: faker.lorem.word(),
        url: faker.internet.url(),
      },
    });
    // Validate status code
    expect(response.status()).toBe(401);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Unauthorized");
  });

  test("Update Feed - Not Found", async ({ request }) => {
    const not_feed_id = faker.string.uuid();
    const response = await request.put(`/v2/feeds/${not_feed_id}`, {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
      data: {
        name: faker.lorem.word(),
        url: faker.internet.url(),
      },
    });
    // Validate status code
    expect(response.status()).toBe(404);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Feed don't exsist");
  });

  test("Update Feed - Invalid URL", async ({ request }) => {
    const response = await request.put(`/v2/feeds/${feed_id}`, {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
      data: {
        name: faker.lorem.word(),
        url: "Invalid URL",
      },
    });
    // Validate status code
    expect(response.status()).toBe(400);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Invalid URL");
  });

  test("Update Feed - URL Exsits", async ({ request }) => {
    const response = await request.put(`/v2/feeds/${feed_id}`, {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
      data: {
        name: faker.lorem.word(),
        url: url2,
      },
    });
    // Validate status code
    expect(response.status()).toBe(409);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Duplicate feed exist");
  });


  test("Update Feed - Feed not owned by user", async ({ request }) => {
    const response = await request.put(`/v2/feeds/${feed_id}`, {
      headers: {
        Authorization: `Bearer ${authToken2}`,
      },
      data: {
        name: faker.lorem.word(),
        url: faker.internet.url(),
      },
    });
    // Validate status code
    expect(response.status()).toBe(403);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Forbidden");
  });
});

test.describe("Delete Feed", () => {
  test("Delete Feed", async ({ request }) => {
    const response = await request.delete(`/v2/feeds/${feed_id}`, {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    });
    // Validate status code
    expect(response.status()).toBe(204);
  });

  test("Delete Feed - Not Authorized", async ({ request }) => {
    const response = await request.delete(`/v2/feeds/${feed_id}`, {});
    // Validate status code
    expect(response.status()).toBe(401);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Unauthorized");
  });

  test("Delete Feed - Not Found", async ({ request }) => {
    const not_feed_id = faker.string.uuid();
    const response = await request.delete(`/v2/feeds/${not_feed_id}`, {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    });
    // Validate status code
    expect(response.status()).toBe(404);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Feed don't exsist");
  });

  test("Delete Feed - Feed not owned by user", async ({ request }) => {
    const response = await request.delete(`/v2/feeds/${feed_id}`, {
      headers: {
        Authorization: `Bearer ${authToken2}`,
      },
    });
    // Validate status code
    expect(response.status()).toBe(403);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Forbidden");
  });
});

test.describe("Get Feeds", () => {
  test("Get all feeds", async ({ request }) => {
    const response = await request.get(`/v2/feeds`, {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    });
    // Validate status code
    expect(response.status()).toBe(200);
    // Validate response body
    const json = await response.json();
    expect(Array.isArray(json)).toBe(true);
    expect(json.length).toBeGreaterThanOrEqual(2);
  });

  test("Get all feeds - Not Authorized", async ({ request }) => {
    const response = await request.get(`/v2/feeds`);
    // Validate status code
    expect(response.status()).toBe(401);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Unauthorized");
  });
});