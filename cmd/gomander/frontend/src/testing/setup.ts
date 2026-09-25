import "@testing-library/jest-dom/vitest";
import { cleanup } from "@testing-library/react";
import { afterEach } from "vitest";

// Testing Library only unmounts on its own when Vitest runs with globals.
afterEach(cleanup);
