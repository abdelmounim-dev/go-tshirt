# 🎯 Project Goal

Build a **clean, maintainable, and well-tested REST API backend** in Go that manages products for an online T-shirt store.
This API should fully support the product lifecycle (CRUD), product options (color, size, stock), recommendations, and shopping cart features.

**Key Technologies:**

*   **Programming Language:** Go
*   **Web Framework:** Gin
*   **ORM:** GORM
*   **Database:** SQLite

**Architecture:**

The project follows a layered architecture:

*   **`cmd/server`:** The main application entry point.
*   **`internal/api`:** Defines the API routes and handlers.
*   **`internal/service`:** Contains the business logic for managing products.
*   **`internal/repository`:** Implements the database operations for products.
*   **`internal/models`:** Defines the `Product` data model.
*   **`internal/config`:** Manages application configuration.
*   **`internal/db`:** Handles the database connection.

# 📋 Planned Features and Tasks

### Phase 1: Core Product CRUD (Already Scaffolded)
- List all products
- Retrieve product by ID
- Create new product
- Update product details
- Delete product

### Phase 2: Validation and Error Handling
- Add input validation for product fields (e.g., price > 0, sizes allowed).
- Return meaningful error messages and HTTP status codes.
- Test validation logic thoroughly.

### Phase 3: Product Options and Availability
- Support selecting colors (`black`, `white`) and sizes (`XS`–`XXL`).
- Grey out unavailable sizes/colors in the frontend API responses.
- Track stock availability per color/size.
- Extend DB schema and models as needed.
- Test business rules for availability.

### Phase 4: Shopping Cart API
- Design Cart resource: adding/removing items, quantities.
- Support session or user association (simplified).
- Ensure cart integrity (e.g., don’t add out-of-stock products).
- CRUD endpoints for cart.
- Full test coverage.

### Phase 5: Recommendations API
- Design “Recommended for you” endpoint.
- Return similar or related products with images and prices.
- Simple logic is fine initially (e.g., same category or color).
- Test response format and data correctness.

### Phase 6: Infrastructure & DevOps (Optional)
- Dockerize the app.
- Add CI pipeline for tests and linting.
- Add database migration tooling.
- Prepare for future PostgreSQL migration.

# Building and Running

**Build:**

```bash
go build ./cmd/server
```

**Run:**

```bash
go run ./cmd/server/main.go
```

The server will start on the address specified in the configuration (default: `:8080`).

**Testing:**

```bash
go test ./...
```

# 🧩 Development Principles

### 1. **Test-Driven Development (TDD) is Mandatory**
- Every feature, bug fix, or refactor must start with writing **failing tests** that define the expected behavior.
- Then implement the minimum code to pass the tests.
- Refactor for readability and performance, keeping tests green.
- Use Go’s standard `testing` package, plus `testify` or similar assertion libs.
- Include both **unit tests** and **integration tests** where applicable.

### 2. **Step-by-Step Feature Development**
- No feature is too small to design and review.
- Each feature should be:
  1. **Explicitly designed and agreed upon** (via PR description, comments, or chat).
  2. **Implemented fully with tests** before moving to the next feature.
  3. **Reviewed and approved** (by a human or AI reviewer).

### 3. **Design First, Code Later**
- Before writing any implementation code:
  - Write down the API contract: endpoints, request/response schemas, error codes.
  - Design data models (Go structs, DB schema).
  - Discuss and finalize logic flows (e.g., how “Add to Cart” behaves with stock).
- Get explicit sign-off from stakeholders or your lead (could be you or the AI supervisor).
- Only then implement.

### 4. **Code Quality and Style**
- Follow idiomatic Go best practices.
- Use consistent naming and folder structures.
- Keep handlers thin; business logic should live in services or domain layers.
- Use interfaces to enable mocking and easier testing.
- Document exported functions and structs.

### 5. **Commit Often and Clearly**
- Every change must be committed separately with a clear, descriptive message.
- Commit in the order of the TDD workflow:
  - First commit the **failing tests** (`test:` prefix).
  - Then commit the **feature implementation** that passes the tests (`feat:` prefix).
  - Finally, commit any **refactoring or cleanup** with tests still passing (`refactor:` prefix).
- This ensures full traceability of every code change and makes the review process smoother.

# 🧑‍💻 Workflow for Each Feature

1. **Design & Agreement**
   - Create or update API specs.
   - Define data models.
   - Document edge cases.

2. **Write Tests**
   - Start with failing tests covering all scenarios.
   - Use table-driven tests where appropriate.
   - Commit: `test: add failing tests for <feature>`

3. **Implement Code**
   - Keep functions small and focused.
   - Handle errors clearly.
   - Integrate with DB via GORM.
   - Commit: `feat: implement <feature> to pass tests`

4. **Run Tests**
   - All tests must pass before pushing.
   - Write additional tests for uncovered cases.

5. **Code Review and Refactor**
   - Clean code for readability or performance.
   - Commit: `refactor: improve <feature> implementation`

6. **Merge & Deploy**
   - After approval, merge to main branch.
   - Deploy manually or via CI pipeline.