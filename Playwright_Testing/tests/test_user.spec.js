import { test, expect } from "@playwright/test";
import { faker } from "@faker-js/faker";
//import { DateTime } from "luxon"; // Luxon is for datetime generation

// Set the BASE URL in playwright.config.js to http://127.0.0.1:8080
let authToken, email, password, email2, authToken2;
email = faker.internet.email();
password = faker.internet.password();
const name = faker.person.firstName();
email2 = faker.internet.email();
test.beforeEach("Credentials for testing", async ({ request }) => {
  // Create User
  const response = await request.post("/v1/user", {
    data: {
      email: email,
      password: password,
      name: name,
    },
  });
  // Create User 2
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
  })
  const json2 = await loginResponse2.json();
  authToken2 = json2.token;
  // Validate status codes
  expect(response.status()).toBe(201);
  expect(loginResponse.status()).toBe(200);

  const json = await loginResponse.json();
  // Get authToken
  authToken = json.token;
});
test.afterEach("Remove Credentials", async ({ request }) => {
  const response = await request.delete("/v1/user", {
    headers: {
      Authorization: `Bearer ${authToken}`,
    },
  });
  const response2 = await request.delete("/v1/user", {
    headers: {
      Authorization: `Bearer ${authToken2}`,
    },
  });

});
test.describe("Create User Tests", () => {
  test("Create user", async ({ request }) => {
    const email_ex = faker.internet.email();
    const password_ex = faker.internet.password();
    const name_ex = faker.person.firstName();

    const response = await request.post("/v1/user", {
      data: {
        email: email_ex,
        password: password_ex,
        name: name_ex,
      },
    });

    const json = await response.json();
    // Validate status code
    expect(response.status()).toBe(201);
    // Validate response body
    expect(json).toHaveProperty("email", email_ex);
    expect(json).toHaveProperty("name", name_ex);
    expect(json).toHaveProperty("id");
    // Login this user
    const loginResponse = await request.post("/v1/login", {
      form: {
        username: email_ex,
        password: password_ex,
      },
    });
    // Validate status code
    expect(loginResponse.status()).toBe(200);
    const loginJson = await loginResponse.json();
    // Get authToken
    const authToken_ex = loginJson.token;

    // Remove this user
    const deleteResponse = await request.delete("/v1/user", {
      headers: {
        Authorization: `Bearer ${authToken_ex}`,
      },
    });
    // Validate status code
    expect(deleteResponse.status()).toBe(204);
  });

  test("Create User - Unvalid Email", async ({ request }) => {
    const response = await request.post("/v1/user", {
      data: {
        email: "test",
        password: password,
        name: name,
      },
    });
    // Validate status code
    expect(response.status()).toBe(400);
    expect(response.ok()).toBeFalsy();
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Invalid email");
  });

  test("Create User - Email exists", async ({ request }) => {
    const response = await request.post("/v1/user", {
      data: {
        email: email,
        password: password,
        name: name,
      },
    });
    // Validate status code
    expect(response.status()).toBe(409);
    expect(response.ok()).toBeFalsy();
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Account already exists");
  });
});
test.describe("Login Tests", () => {
  test("Login", async ({ request }) => {
    const response = await request.post("/v1/login", {
      form: {
        username: email,
        password: password,
      },
    });

    const json = await response.json();
    // Validate status code
    expect(response.status()).toBe(200);
    // Validate response body
    expect(json).toHaveProperty("token");
    expect(json).toHaveProperty("token_type");
    authToken = json.token;
  });

  test("Login - Not valid Username", async ({ request }) => {
    const response = await request.post("/v1/login", {
      form: {
        username: "invalid",
        password: password,
      },
    });
    // Validate status code
    expect(response.status()).toBe(400);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Invalid email");
  });

  test("Login - Not valid Password", async ({ request }) => {
    const response = await request.post("/v1/login", {
      form: {
        username: email,
        password: "invalid",
      },
    });
    // Validate status code
    expect(response.status()).toBe(401);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Wrong password");
  });

  test("Login - Not valid user", async ({ request }) => {
    const emai_ex = faker.internet.email();
    const response = await request.post("/v1/login", {
      form: {
        username: emai_ex,
        password: "invalid",
      },
    });
    // Validate status code
    expect(response.status()).toBe(404);
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "User not found");
  });
});
test.describe("Update user", () => {
  test("Update user", async ({ request }) => {
    const new_email = faker.internet.email();
    const new_name = faker.person.firstName();

    const response = await request.put("/v1/user", {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
      data: {
        email: new_email,
        name: new_name,
      },
    });
    // Validate status code
    expect(response.status()).toBe(200);
    expect(response.ok()).toBeTruthy();
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("email", new_email);
    expect(json).toHaveProperty("name", new_name);
  });

  test("Update user - Not Authorized", async ({ request }) => {
    const response = await request.put("/v1/user", {
      data: {
        email: "new@gmail.com",
        name: name,
      },
    });
    // Validate status code
    expect(response.status()).toBe(401);
    expect(response.ok()).toBeFalsy();
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Unauthorized");
  });

  test("Update user - Not valid email", async ({ request }) => {
    const response = await request.put("/v1/user", {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
      data: {
        email: "test",
        name: name,
      },
    });
    // Validate status code
    expect(response.status()).toBe(400);
    expect(response.ok()).toBeFalsy();
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Invalid email");
  });

  test("Update user - Email exists", async ({ request }) => {
    const response = await request.put("/v1/user", {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
      data: {
        email: email2,
        name: name,
      },
    });
    // Validate status code
    expect(response.status()).toBe(409);
    expect(response.ok()).toBeFalsy();
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Account already exists");
  });
});
test.describe("Delete user", () => {
  test("Delete user", async ({ request }) => {
    const response = await request.delete("/v1/user", {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    });
    // Validate status code
    expect(response.status()).toBe(204);
    expect(response.ok()).toBeTruthy();
  });

  test("Delete user - Not Authorized", async ({ request }) => {
    const response = await request.delete("/v1/user");
    // Validate status code
    expect(response.status()).toBe(401);
    expect(response.ok()).toBeFalsy();
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Unauthorized");
  });

  test("Delete user - Not Found", async ({ request }) => {
    await request.delete("/v1/user", {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    });

    const response = await request.delete("/v1/user", {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    });
    // Validate status code
    expect(response.status()).toBe(404);
    expect(response.ok()).toBeFalsy();
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "User don't exsist");
  });
});
test.describe("Get user", () => {
  test("Get a user", async ({ request }) => {
    const response = await request.get("/v1/user", {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    });
    // Validate status code
    expect(response.status()).toBe(200);
    expect(response.ok()).toBeTruthy();
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("email", email);
    expect(json).toHaveProperty("name", name);
  });

  test("Get user - Not Authorized", async ({ request }) => {
    const response = await request.get("/v1/user");
    // Validate status code
    expect(response.status()).toBe(401);
    expect(response.ok()).toBeFalsy();
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "Unauthorized");
  });

  test("Get user - Not Found", async ({ request }) => {
    // Delete the user prior
    await request.delete("/v1/user", {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    });
    // Re-access the user
    const response = await request.get("/v1/user", {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    });
    // Validate status code
    expect(response.status()).toBe(404);
    expect(response.ok()).toBeFalsy();
    // Validate response body
    const json = await response.json();
    expect(json).toHaveProperty("error", "User don't exsist");
  });
});
