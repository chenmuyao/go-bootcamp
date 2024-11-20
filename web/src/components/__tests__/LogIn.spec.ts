import { describe, it, expect, beforeEach, vi } from 'vitest'

import { mount, VueWrapper } from '@vue/test-utils'
import LogIn from '../LogIn.vue'

describe('Login', () => {
  let wrapper: VueWrapper;

  beforeEach(() => {
    wrapper = mount(LogIn)
  });

  it('renders the form with inputs and a button', () => {
    expect(wrapper.find('input[type="email"]').exists()).toBe(true);
    expect(wrapper.find('input[type="password"]').exists()).toBe(true);
    expect(wrapper.find('button').exists()).toBe(true);
  });

  it('updates the email input value correctly', async () => {
    const emailInput = wrapper.get('input[type="email"]');
    await emailInput.setValue('test@example.com');
    expect(emailInput.element.value).toBe('test@example.com');
  });

  it('disable the login button when inputs are invalid', async () => {
    const loginButton = wrapper.get('button');
    expect(loginButton.attributes('disabled')).toBeDefined;
  });

  it('enables the login button when inputs are valid', async () => {
    await wrapper.get('input[type="email"]').setValue('test@example.com');
    await wrapper.get('input[type="password"]').setValue('password123!');

    const loginButton = wrapper.get('button');
    expect(loginButton.attributes('disabled')).toBeUndefined();
  });

  it('shows error if the email is not a valid email', async () => {
    await wrapper.get('input[type="email"]').setValue('not a email');
    const error = wrapper.find('.error-message');

    expect(error.exists()).toBe(true);
    expect(error.text()).toBe('Please enter a valid email');
  })

  it('emits a login event with correct data on valid form submission', async () => {
    const handleLoginSpy = vi.spyOn(wrapper.vm.$, 'emit');
    const formData = {
      email: 'test@example.com',
      password: 'password123!',
    };

    await wrapper.get('input[type="email"]').setValue('test@example.com');
    await wrapper.get('input[type="password"]').setValue('password123!');

    await wrapper.get('form').trigger('submit');

    expect(handleLoginSpy).toHaveBeenCalledWith('login', {
      email: formData.email,
      password: formData.password,
    });
  });
});


