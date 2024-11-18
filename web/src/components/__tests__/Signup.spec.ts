import { describe, it, expect, beforeEach, vi } from 'vitest'

import { mount, VueWrapper } from '@vue/test-utils'
import SignUp from '../SignUp.vue'

interface SignupFormData {
  email: string;
  password: string;
  confirmPassword: string;
}

describe('Signup', () => {
  let wrapper: VueWrapper;

  beforeEach(() => {
    wrapper = mount(SignUp)
  });

  it('renders the form with inputs and a button', () => {
    expect(wrapper.find('input[type="#email"]').exists()).toBe(true);
    expect(wrapper.find('input[type="#password"]').exists()).toBe(true);
    expect(wrapper.find('input[type="#confirm-password"]').exists()).toBe(true);
    expect(wrapper.find('button"]').exists()).toBe(true);
  });

  it('updates the email input value correctly', async () => {
    const emailInput = wrapper.get('input[type="email"]');
    await emailInput.setValue('test@example.com');
    expect(emailInput.element.nodeValue).toBe('test@example.com');
  });

  it('disable the signup button when inputs are invalid', async () => {
    const signupButton = wrapper.get('button');
    expect(signupButton.attributes('disabled')).toBe('true');
  });

  it('enables the signup button when inputs are valid', async () => {
    await wrapper.get('input[type="email"]').setValue('test@example.com');
    await wrapper.get('input[type="password"]').setValue('password123!');
    await wrapper.get('input#confirm-password').setValue('password123!');

    const signupButton = wrapper.get('button');
    expect(signupButton.attributes('disabled')).toBeUndefined();
  });

  it('shows error if passwords do not match', async () => {
    await wrapper.get('input[type="password"]').setValue('password123!');
    await wrapper.get('input#confirm-password').setValue('notmatch');
    const error = wrapper.find('.error-message'); // TODO: error message class

    expect(error.exists()).toBe(true);
    expect(error.text()).toBe('Passwords do not match');
  });

  it('shows error if passwords do not have letters, numbers, at least one special characters', async () => {
    await wrapper.get('input[type="email"]').setValue('test@example.com');
    await wrapper.get('input[type="password"]').setValue('password123');
    await wrapper.get('input#confirm-password').setValue('password123');

    await wrapper.get('button').trigger('click');

    const error = wrapper.find('.error-message'); // TODO: error message class

    expect(error.exists()).toBe(true);
    expect(error.text()).toBe('Password should at least have letters, numbers and at least one special characters');
  });

  it('shows error if the email is not a valid email', async () => {
    await wrapper.get('input[type="email"]').setValue('not a email');
    const error = wrapper.find('.error-message');

    expect(error.exists()).toBe(true);
    expect(wrapper.text()).toBe('Input is not a valid email');
  })

  it('emits a signup event with correct data on valid form submission', async () => {
    const mockSignup = vi.fn()
    const formData: SignupFormData = {
      email: 'test@example.com',
      password: 'password123!',
      confirmPassword: 'password123!',
    };

    wrapper.vm.$emit = mockSignup; // Mock the emit method

    await wrapper.get('input[type="email"]').setValue('test@example.com');
    await wrapper.get('input[type="password"]').setValue('password123');
    await wrapper.get('input#confirm-password').setValue('password123');

    await wrapper.get('button').trigger('click');

    expect(mockSignup).toHaveBeenCalledWith('signup', {
      email: formData.email,
      password: formData.password,
      confirmPassword: formData.confirmPassword,
    });
  });
});

