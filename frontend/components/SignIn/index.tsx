import Link from "next/link";
import React, { useCallback, useContext, useMemo } from "react";
import { Alert, Button, Card, Container, Form } from "react-bootstrap";
import { CanonicalLink, Description, Layout, Title } from "..";
import { AuthContext } from "../../contexts";
import { api, base64Encode, getIsHosted, setUserAccessToken } from "../../helpers";
import { useForm } from "../../hooks";

type Form = {
  email: string;
  password: string;
};

export function SignIn() {
  const { mutate } = useContext(AuthContext);
  const [values, , updateValue, errors, setErrors] = useForm<Form>({ email: "", password: "" });

  const errorNode = useMemo(() => (errors.detail !== undefined ? (
    <Alert variant="danger">
      {errors.detail}
    </Alert>
  ) : null), [errors.detail]);

  const handleSubmit: React.FormEventHandler<HTMLFormElement> = useCallback(async (event) => {
    event.preventDefault();

    const response = await api.post("/user-access-tokens", {}, {
      headers: {
        authorization: `basic ${base64Encode(`${values.email}:${values.password}`)}`,
      },
    });
    const responseJson = await response.json();

    if (response.ok) {
      setUserAccessToken(responseJson.token);
      await mutate();
    } else {
      setErrors(responseJson);
    }
  }, [mutate, setErrors, values.email, values.password]);

  return (
    <Layout kind="app">
      <CanonicalLink path="/sign-in" />

      <Title>Sign in</Title>

      <Description>Sign in to your PoeticMetric account and gain access to powerful website analytics features.</Description>

      <Container className="py-5">
        <div className="text-center">
          <h1>Sign in to continue</h1>

          {getIsHosted() ? (
            <div className="mt-3">
              {"Don't have an account? "}

              <Link href="/sign-up">Sign up</Link>
            </div>
          ) : null}
        </div>

        <Card className="mt-4 mx-auto mw-32rem">
          <Card.Body>
            {errorNode}

            <Form onSubmit={handleSubmit}>
              <Form.Group controlId="email">
                <Form.Label>E-mail address</Form.Label>

                <Form.Control
                  isInvalid={errors.email !== undefined}
                  name="email"
                  onChange={updateValue}
                  required
                  type="email"
                  value={values.email}
                />
              </Form.Group>

              <Form.Group className="mt-2" controlId="password">
                <Form.Label>Password</Form.Label>

                <Form.Control
                  isInvalid={errors.password !== undefined}
                  name="password"
                  onChange={updateValue}
                  required
                  type="password"
                  value={values.password}
                />
              </Form.Group>

              <div className="d-grid mt-4">
                <Button type="submit" variant="primary">Sign in</Button>
              </div>

              <div className="align-items-center d-flex flex-column mt-2">
                <Link
                  className="fs-sm fw-semibold"
                  href={{ pathname: "/password-recovery", query: values.email === "" ? undefined : { email: values.email } }}
                >
                  Forgot password?
                </Link>
              </div>
            </Form>
          </Card.Body>
        </Card>
      </Container>
    </Layout>
  );
}
