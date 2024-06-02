# Money Tracker API Documentation

## Server

http://localhost:8080

## Authentication Method

Access Token expires in 5 minutes

Refresh token expires in 30 day

```json
{
  "headers": {
    "Authorization": "Bearer <token>"
  }
}
```

## Endpoints

### Authentication

#### POST /api/v1/google/callback

- Request

```ts
type Body = {
  // Google Auth Code
  code: string;
};
```

- Response

```ts
type Response = {
  code: 201;
  data: {
    access_token: string;
    refresh_token: string;
    access_token_expires_in: number;
    refresh_token_expires_in: number;
  };
};
```

#### POST /api/v1/refresh

- Auth Required
- Headers

  ```json
  {
    "headers": {
      "Authorization": "Bearer <refresh>"
    }
  }
  ```

- Response

```ts
type Response = {
  code: 201;
  data: {
    access_token: string;
    refresh_token: string;
    access_token_expires_in: number;
    refresh_token_expires_in: number;
  };
};
```

#### POST /api/v1/logout

- Auth Required
- Response

```ts
type Response = {
  code: 200;
  data: 'LOGOUT_SUCCESS';
};
```

#### Authentication Known Errors

- Empty JWT Token

  ```ts
  type Response = {
    code: 401;
    data: 'UNAUTHORIZED';
  };
  ```

- JWT Signature Invalid / Token Invalid

  ```ts
  type Response = {
    code: 401;
    data: 'TOKEN_INVALID';
  };
  ```

  - JWT Expired

  ```ts
  type Response = {
    code: 401;
    data: 'TOKEN_EXPIRED';
  };
  ```

### Category

#### POST /api/v1/category

- Auth Required
- Request

```ts
type Body = {
  name: string;
  type: 'income' | 'expense';
};
```

- Response

```ts
type Response = {
  code: 201;
  data: {
    id: string;
    name: string;
    slug: string;
    type: 'income' | 'expense';
  };
};
```

#### POST /api/v1/category/subcategory

- Auth Required
- Request

```ts
type Body = {
  name: string;
  category_id: string;
};
```

- Response

```ts
type Response = {
  code: 201;
  data: {
    id: string;
    name: string;
    slug: string;
    category_id: string;
  };
};
```

#### PUT /api/v1/category/subcategory/:id

- Auth Required
- Parameters

```ts
type Params = {
  id: number;
};
```

- Request

```ts
type Body = {
  name: string;
};
```

- Response

```ts
type Response = {
  code: 200;
  data: {
    id: string;
    name: string;
    slug: string;
    category_id: string;
  };
};
```

#### DELETE /api/v1/category/subcategory/:id

- Auth Required
- Parameters

```ts
type Params = {
  id: number;
};
```

- Response

```ts
type Response = {
  code: 200;
  data: 'SUBCATEGORY_DELETED';
};
```

### Transaction

#### POST /api/v1/transaction

- Auth Required
- Request

  ```ts
  type Body = {
    amount: number;
    category_id: string;
    subcategory_id?: string;
    transaction_at: string; // YYYY-MM-DD
    description?: string;
  };
  ```

- Response

  ```ts
  type Response = {
    code: 201;
    data: {
      id: string;
      amount: number;
      category: {
        id: string;
        name: string;
        slug: string;
        type: 'income' | 'expense';
      };
      subcategory: {
        id: string;
        name: string;
        slug: string;
      } | null;
      transaction_at: string;
      description: string | null;
      created_at: string;
      updated_at: string;
    };
  };
  ```

#### PUT /api/v1/transaction/:id

- Auth Required
- Parameters

  ```ts
  type Params = {
    id: number;
  };
  ```

- Request

  ```ts
  type Body = {
    amount: number;
    category_id: string;
    subcategory_id?: string;
    transaction_at: string; // YYYY-MM-DD
    description?: string;
  };
  ```

- Response

  ```ts
  type Response = {
    code: 201;
    data: {
      id: string;
      amount: number;
      category: {
        id: string;
        name: string;
        slug: string;
        type: 'income' | 'expense';
      };
      subcategory: {
        id: string;
        name: string;
        slug: string;
      } | null;
      transaction_at: string;
      description: string | null;
      created_at: string;
      updated_at: string;
    };
  };
  ```

#### DELETE /api/v1/transaction/:id

- Auth Required
- Parameters

  ```ts
  type Params = {
    id: number;
  };
  ```

- Response

  ```ts
  type Response = {
    code: 200;
    data: null;
  };
  ```

#### GET /api/v1/transaction/:id

- Auth Required
- Parameters

  ```ts
  type Params = {
    id: number;
  };
  ```

- Response

  ```ts
  type Response = {
    code: 200;
    data: {
      id: string;
      amount: number;
      category: {
        id: string;
        name: string;
        slug: string;
        type: 'income' | 'expense';
      };
      subcategory: {
        id: string;
        name: string;
        slug: string;
      } | null;
      transaction_at: string;
      description: string | null;
      created_at: string;
      updated_at: string;
    };
  };
  ```

#### GET /api/v1/transaction

- Auth Required
- Query Parameters

  ```ts
  type Query = {
    page: number; // min: 1,
    limit: number; // min: 1, max: 20
    category_id?: string;
    subcategory_id?: string;
    started_at?: string; // YYYY-MM-DD
    ended_at?: string; // YYYY-MM-DD
    type: 'income' | 'expense';
  };
  ```

- Response

  ```ts
  type Response = {
    code: 200;
    data: Array<{
      id: string;
      amount: number;
      category: {
        id: string;
        name: string;
        slug: string;
        type: 'income' | 'expense';
      };
      subcategory: {
        id: string;
        name: string;
        slug: string;
      } | null;
      transaction_at: string;
      description: string | null;
      created_at: string;
      updated_at: string;
    }> | null;
    metadata: {
      page: number;
      limit: number;
      total_docs: number;
      total_pages: number;
      has_next_page: boolean;
    };
  };
  ```
