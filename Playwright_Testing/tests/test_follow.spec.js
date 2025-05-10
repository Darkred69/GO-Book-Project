import { test, expect } from "@playwright/test";
import { faker } from "@faker-js/faker";

let email, password, authToken, feed_id, user_id, feed_id2;
test.beforeAll("Credentials - User", async ({ request }) => {
  email = faker.internet.email();
  password = faker.internet.password();
  const name = faker.person.firstName();
  // Create User
  const response = await request.post("/v1/user", {
    data: {
      email: email,
      password: password,
      name: name,
    },
  });
  const body = await response.json();
  user_id = body.id;
  // Login User
  const loginResponse = await request.post("/v1/login", {
    form: {
      username: email,
      password: password,
    },
  });

  const json = await loginResponse.json();
  authToken = json.token;

  // Validate status code
  expect(loginResponse.status()).toBe(200);
  expect(response.status()).toBe(201);
});

test.beforeAll("Credentials - Feed", async ({ request }) => {
  // Feed 1
  const title = faker.lorem.word();
  const url = faker.internet.url();
  const response = await request.post("/v2/feeds", {
    headers: {
      Authorization: `Bearer ${authToken}`,
    },
    data: {
      title: title,
      url: url,
    },
  });

  const json = await response.json();
  feed_id = json.id;

  // Feed 2
  const title2 = faker.lorem.word();
  const url2 = faker.internet.url();
  const response2 = await request.post("/v2/feeds", {
    headers: {
      Authorization: `Bearer ${authToken}`,
    },
    data: {
      title: title2,
      url: url2,
    },
  });
  const json2 = await response2.json();
  feed_id2 = json2.id;

  // Validate status code
  expect(response.status()).toBe(201);
});

test.beforeAll("Credentials - Follow", async ({ request }) => {
  const response = await request.post(`/v3/follow`, {
    headers: {
      Authorization: `Bearer ${authToken}`,
    },
    data: {
      feed_id: feed_id2,
    },
  });

  // Validate status code
  expect(response.status()).toBe(201);
});

test.afterAll("Remove Credentials", async ({ request }) => {
  await request.delete(`v2/feeds/${feed_id}`, {
    headers: {
      Authorization: `Bearer ${authToken}`,
    },
  });

  await request.delete(`v2/feeds/${feed_id2}`, {
    headers: {
      Authorization: `Bearer ${authToken}`,
    },
  });

  await request.delete("/v1/user", {
    headers: {
      Authorization: `Bearer ${authToken}`,
    },
  });
});
test.describe("Unfollow Feed", () => {
  test("Unfollow Feed", async ({ request }) => {
    const response = await request.delete(`/v3/follow/${feed_id2}`, {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    });

    // Validate status code
    expect(response.status()).toBe(204);
  });

  test("Unfollow Feed - Feed Not Found", async ({ request }) => {
    const not_feed_id = faker.string.uuid();
    const response = await request.delete(`/v3/follow/${not_feed_id}`, {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    });

    // Validate status code
    expect(response.status()).toBe(404);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Feed not found");
  });

  test("Unfollow Feed - Never followed", async ({ request }) => {
    await request.delete(`/v3/follow/${feed_id2}`, {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    });
    const response = await request.delete(`/v3/follow/${feed_id2}`, {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    });
    // Validate status code
    expect(response.status()).toBe(404);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Feed not followed");
  });
});

test.describe("Follow Feed", () => {
  test("Follow Feed", async ({ request }) => {
    const response = await request.post(`/v3/follow`, {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
      data: {
        feed_id: feed_id,
      },
    });

    // Validate status code
    expect(response.status()).toBe(201);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("feed_id", feed_id);
    expect(json).toHaveProperty("user_id", user_id);
  });

  test("Follow Feed - Feed Not Found", async ({ request }) => {
    const not_feed_id = faker.string.uuid();
    const response = await request.post(`/v3/follow`, {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
      data: {
        feed_id: not_feed_id,
      },
    });

    // Validate status code
    expect(response.status()).toBe(404);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Feed not found");
  });

  test("Follow Feed - Feed Already Followed", async ({ request }) => {
    const response = await request.post(`/v3/follow`, {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
      data: {
        feed_id: feed_id2,
      },
    });

    // Validate status code
    expect(response.status()).toBe(409);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Feed already followed");
  });

  test("Follow Feed - Not Authorized", async ({ request }) => {
    const response = await request.post(`/v3/follow`, {
      data: {
        feed_id: feed_id2,
      },
    });
    // Validate status code
    expect(response.status()).toBe(401);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Unauthorized");
  });
});

test.describe("Get follows", () => {
  test("Get all follows", async ({ request }) => {
    const response = await request.get(`/v3/follow`, {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    });

    // Validate status code
    expect(response.status()).toBe(200);
    // Validate response body
    const json = await response.json();
    expect(Array.isArray(json)).toBe(true);
    expect(json.length).toBe(1);
  });

  test("Get all follows - Not Authorized", async ({ request }) => {
    const response = await request.get(`/v3/follow`);
    // Validate status code
    expect(response.status()).toBe(401);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Unauthorized");
  });
});
