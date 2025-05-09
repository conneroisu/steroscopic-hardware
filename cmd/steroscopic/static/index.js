var F0 = (function () {
    let htmx = {
      onLoad: null,
      process: null,
      on: null,
      off: null,
      trigger: null,
      ajax: null,
      find: null,
      findAll: null,
      closest: null,
      values: function (J, Y) {
        return getInputValues(J, Y || "post").values;
      },
      remove: null,
      addClass: null,
      removeClass: null,
      toggleClass: null,
      takeClass: null,
      swap: null,
      defineExtension: null,
      removeExtension: null,
      logAll: null,
      logNone: null,
      logger: null,
      config: {
        historyEnabled: !0,
        historyCacheSize: 10,
        refreshOnHistoryMiss: !1,
        defaultSwapStyle: "innerHTML",
        defaultSwapDelay: 0,
        defaultSettleDelay: 20,
        includeIndicatorStyles: !0,
        indicatorClass: "htmx-indicator",
        requestClass: "htmx-request",
        addedClass: "htmx-added",
        settlingClass: "htmx-settling",
        swappingClass: "htmx-swapping",
        allowEval: !0,
        allowScriptTags: !0,
        inlineScriptNonce: "",
        inlineStyleNonce: "",
        attributesToSettle: ["class", "style", "width", "height"],
        withCredentials: !1,
        timeout: 0,
        wsReconnectDelay: "full-jitter",
        wsBinaryType: "blob",
        disableSelector: "[hx-disable], [data-hx-disable]",
        scrollBehavior: "instant",
        defaultFocusScroll: !1,
        getCacheBusterParam: !1,
        globalViewTransitions: !1,
        methodsThatUseUrlParams: ["get", "delete"],
        selfRequestsOnly: !0,
        ignoreTitle: !1,
        scrollIntoViewOnBoost: !0,
        triggerSpecsCache: null,
        disableInheritance: !1,
        responseHandling: [
          { code: "204", swap: !1 },
          { code: "[23]..", swap: !0 },
          { code: "[45]..", swap: !1, error: !0 },
        ],
        allowNestedOobSwaps: !0,
      },
      parseInterval: null,
      _: null,
      version: "2.0.4",
    };
    (htmx.onLoad = onLoadHelper),
      (htmx.process = processNode),
      (htmx.on = addEventListenerImpl),
      (htmx.off = removeEventListenerImpl),
      (htmx.trigger = triggerEvent),
      (htmx.ajax = ajaxHelper),
      (htmx.find = find),
      (htmx.findAll = findAll),
      (htmx.closest = closest),
      (htmx.remove = removeElement),
      (htmx.addClass = addClassToElement),
      (htmx.removeClass = removeClassFromElement),
      (htmx.toggleClass = toggleClassOnElement),
      (htmx.takeClass = takeClassForElement),
      (htmx.swap = swap),
      (htmx.defineExtension = defineExtension),
      (htmx.removeExtension = removeExtension),
      (htmx.logAll = logAll),
      (htmx.logNone = logNone),
      (htmx.parseInterval = parseInterval),
      (htmx._ = internalEval);
    let internalAPI = {
        addTriggerHandler,
        bodyContains,
        canAccessLocalStorage,
        findThisElement,
        filterValues,
        swap,
        hasAttribute,
        getAttributeValue,
        getClosestAttributeValue,
        getClosestMatch,
        getExpressionVars,
        getHeaders,
        getInputValues,
        getInternalData,
        getSwapSpecification,
        getTriggerSpecs,
        getTarget,
        makeFragment,
        mergeObjects,
        makeSettleInfo,
        oobSwap,
        querySelectorExt,
        settleImmediately,
        shouldCancel,
        triggerEvent,
        triggerErrorEvent,
        withExtensions,
      },
      VERBS = ["get", "post", "put", "delete", "patch"],
      VERB_SELECTOR = VERBS.map(function (J) {
        return "[hx-" + J + "], [data-hx-" + J + "]";
      }).join(", ");
    function parseInterval(J) {
      if (J == null) return;
      let Y = NaN;
      if (J.slice(-2) == "ms") Y = parseFloat(J.slice(0, -2));
      else if (J.slice(-1) == "s") Y = parseFloat(J.slice(0, -1)) * 1000;
      else if (J.slice(-1) == "m") Y = parseFloat(J.slice(0, -1)) * 1000 * 60;
      else Y = parseFloat(J);
      return isNaN(Y) ? void 0 : Y;
    }
    function getRawAttribute(J, Y) {
      return J instanceof Element && J.getAttribute(Y);
    }
    function hasAttribute(J, Y) {
      return (
        !!J.hasAttribute && (J.hasAttribute(Y) || J.hasAttribute("data-" + Y))
      );
    }
    function getAttributeValue(J, Y) {
      return getRawAttribute(J, Y) || getRawAttribute(J, "data-" + Y);
    }
    function parentElt(J) {
      let Y = J.parentElement;
      if (!Y && J.parentNode instanceof ShadowRoot) return J.parentNode;
      return Y;
    }
    function getDocument() {
      return document;
    }
    function getRootNode(J, Y) {
      return J.getRootNode ? J.getRootNode({ composed: Y }) : getDocument();
    }
    function getClosestMatch(J, Y) {
      while (J && !Y(J)) J = parentElt(J);
      return J || null;
    }
    function getAttributeValueWithDisinheritance(J, Y, Z) {
      let $ = getAttributeValue(Y, Z),
        W = getAttributeValue(Y, "hx-disinherit");
      var X = getAttributeValue(Y, "hx-inherit");
      if (J !== Y) {
        if (htmx.config.disableInheritance)
          if (X && (X === "*" || X.split(" ").indexOf(Z) >= 0)) return $;
          else return null;
        if (W && (W === "*" || W.split(" ").indexOf(Z) >= 0)) return "unset";
      }
      return $;
    }
    function getClosestAttributeValue(J, Y) {
      let Z = null;
      if (
        (getClosestMatch(J, function ($) {
          return !!(Z = getAttributeValueWithDisinheritance(
            J,
            asElement($),
            Y,
          ));
        }),
        Z !== "unset")
      )
        return Z;
    }
    function matches(J, Y) {
      let Z =
        J instanceof Element &&
        (J.matches ||
          J.matchesSelector ||
          J.msMatchesSelector ||
          J.mozMatchesSelector ||
          J.webkitMatchesSelector ||
          J.oMatchesSelector);
      return !!Z && Z.call(J, Y);
    }
    function getStartTag(J) {
      let Z = /<([a-z][^\/\0>\x20\t\r\n\f]*)/i.exec(J);
      if (Z) return Z[1].toLowerCase();
      else return "";
    }
    function parseHTML(J) {
      return new DOMParser().parseFromString(J, "text/html");
    }
    function takeChildrenFor(J, Y) {
      while (Y.childNodes.length > 0) J.append(Y.childNodes[0]);
    }
    function duplicateScript(J) {
      let Y = getDocument().createElement("script");
      if (
        (forEach(J.attributes, function (Z) {
          Y.setAttribute(Z.name, Z.value);
        }),
        (Y.textContent = J.textContent),
        (Y.async = !1),
        htmx.config.inlineScriptNonce)
      )
        Y.nonce = htmx.config.inlineScriptNonce;
      return Y;
    }
    function isJavaScriptScriptNode(J) {
      return (
        J.matches("script") &&
        (J.type === "text/javascript" || J.type === "module" || J.type === "")
      );
    }
    function normalizeScriptTags(J) {
      Array.from(J.querySelectorAll("script")).forEach((Y) => {
        if (isJavaScriptScriptNode(Y)) {
          let Z = duplicateScript(Y),
            $ = Y.parentNode;
          try {
            $.insertBefore(Z, Y);
          } catch (W) {
            logError(W);
          } finally {
            Y.remove();
          }
        }
      });
    }
    function makeFragment(J) {
      let Y = J.replace(/<head(\s[^>]*)?>[\s\S]*?<\/head>/i, ""),
        Z = getStartTag(Y),
        $;
      if (Z === "html") {
        $ = new DocumentFragment();
        let X = parseHTML(J);
        takeChildrenFor($, X.body), ($.title = X.title);
      } else if (Z === "body") {
        $ = new DocumentFragment();
        let X = parseHTML(Y);
        takeChildrenFor($, X.body), ($.title = X.title);
      } else {
        let X = parseHTML(
          '<body><template class="internal-htmx-wrapper">' +
            Y +
            "</template></body>",
        );
        ($ = X.querySelector("template").content), ($.title = X.title);
        var W = $.querySelector("title");
        if (W && W.parentNode === $) W.remove(), ($.title = W.innerText);
      }
      if ($)
        if (htmx.config.allowScriptTags) normalizeScriptTags($);
        else $.querySelectorAll("script").forEach((X) => X.remove());
      return $;
    }
    function maybeCall(J) {
      if (J) J();
    }
    function isType(J, Y) {
      return Object.prototype.toString.call(J) === "[object " + Y + "]";
    }
    function isFunction(J) {
      return typeof J === "function";
    }
    function isRawObject(J) {
      return isType(J, "Object");
    }
    function getInternalData(J) {
      let Z = J["htmx-internal-data"];
      if (!Z) Z = J["htmx-internal-data"] = {};
      return Z;
    }
    function toArray(J) {
      let Y = [];
      if (J) for (let Z = 0; Z < J.length; Z++) Y.push(J[Z]);
      return Y;
    }
    function forEach(J, Y) {
      if (J) for (let Z = 0; Z < J.length; Z++) Y(J[Z]);
    }
    function isScrolledIntoView(J) {
      let Y = J.getBoundingClientRect(),
        Z = Y.top,
        $ = Y.bottom;
      return Z < window.innerHeight && $ >= 0;
    }
    function bodyContains(J) {
      return J.getRootNode({ composed: !0 }) === document;
    }
    function splitOnWhitespace(J) {
      return J.trim().split(/\s+/);
    }
    function mergeObjects(J, Y) {
      for (let Z in Y) if (Y.hasOwnProperty(Z)) J[Z] = Y[Z];
      return J;
    }
    function parseJSON(J) {
      try {
        return JSON.parse(J);
      } catch (Y) {
        return logError(Y), null;
      }
    }
    function canAccessLocalStorage() {
      try {
        return (
          localStorage.setItem(
            "htmx:localStorageTest",
            "htmx:localStorageTest",
          ),
          localStorage.removeItem("htmx:localStorageTest"),
          !0
        );
      } catch (Y) {
        return !1;
      }
    }
    function normalizePath(J) {
      try {
        let Y = new URL(J);
        if (Y) J = Y.pathname + Y.search;
        if (!/^\/$/.test(J)) J = J.replace(/\/+$/, "");
        return J;
      } catch (Y) {
        return J;
      }
    }
    function internalEval(str) {
      return maybeEval(getDocument().body, function () {
        return eval(str);
      });
    }
    function onLoadHelper(J) {
      return htmx.on("htmx:load", function (Z) {
        J(Z.detail.elt);
      });
    }
    function logAll() {
      htmx.logger = function (J, Y, Z) {
        if (console) console.log(Y, J, Z);
      };
    }
    function logNone() {
      htmx.logger = null;
    }
    function find(J, Y) {
      if (typeof J !== "string") return J.querySelector(Y);
      else return find(getDocument(), J);
    }
    function findAll(J, Y) {
      if (typeof J !== "string") return J.querySelectorAll(Y);
      else return findAll(getDocument(), J);
    }
    function getWindow() {
      return window;
    }
    function removeElement(J, Y) {
      if (((J = resolveTarget(J)), Y))
        getWindow().setTimeout(function () {
          removeElement(J), (J = null);
        }, Y);
      else parentElt(J).removeChild(J);
    }
    function asElement(J) {
      return J instanceof Element ? J : null;
    }
    function asHtmlElement(J) {
      return J instanceof HTMLElement ? J : null;
    }
    function asString(J) {
      return typeof J === "string" ? J : null;
    }
    function asParentNode(J) {
      return J instanceof Element ||
        J instanceof Document ||
        J instanceof DocumentFragment
        ? J
        : null;
    }
    function addClassToElement(J, Y, Z) {
      if (((J = asElement(resolveTarget(J))), !J)) return;
      if (Z)
        getWindow().setTimeout(function () {
          addClassToElement(J, Y), (J = null);
        }, Z);
      else J.classList && J.classList.add(Y);
    }
    function removeClassFromElement(J, Y, Z) {
      let $ = asElement(resolveTarget(J));
      if (!$) return;
      if (Z)
        getWindow().setTimeout(function () {
          removeClassFromElement($, Y), ($ = null);
        }, Z);
      else if ($.classList) {
        if (($.classList.remove(Y), $.classList.length === 0))
          $.removeAttribute("class");
      }
    }
    function toggleClassOnElement(J, Y) {
      (J = resolveTarget(J)), J.classList.toggle(Y);
    }
    function takeClassForElement(J, Y) {
      (J = resolveTarget(J)),
        forEach(J.parentElement.children, function (Z) {
          removeClassFromElement(Z, Y);
        }),
        addClassToElement(asElement(J), Y);
    }
    function closest(J, Y) {
      if (((J = asElement(resolveTarget(J))), J && J.closest))
        return J.closest(Y);
      else {
        do if (J == null || matches(J, Y)) return J;
        while ((J = J && asElement(parentElt(J))));
        return null;
      }
    }
    function startsWith(J, Y) {
      return J.substring(0, Y.length) === Y;
    }
    function endsWith(J, Y) {
      return J.substring(J.length - Y.length) === Y;
    }
    function normalizeSelector(J) {
      let Y = J.trim();
      if (startsWith(Y, "<") && endsWith(Y, "/>"))
        return Y.substring(1, Y.length - 2);
      else return Y;
    }
    function querySelectorAllExt(J, Y, Z) {
      if (Y.indexOf("global ") === 0)
        return querySelectorAllExt(J, Y.slice(7), !0);
      J = resolveTarget(J);
      let $ = [];
      {
        let Q = 0,
          G = 0;
        for (let B = 0; B < Y.length; B++) {
          let U = Y[B];
          if (U === "," && Q === 0) {
            $.push(Y.substring(G, B)), (G = B + 1);
            continue;
          }
          if (U === "<") Q++;
          else if (U === "/" && B < Y.length - 1 && Y[B + 1] === ">") Q--;
        }
        if (G < Y.length) $.push(Y.substring(G));
      }
      let W = [],
        X = [];
      while ($.length > 0) {
        let Q = normalizeSelector($.shift()),
          G;
        if (Q.indexOf("closest ") === 0)
          G = closest(asElement(J), normalizeSelector(Q.substr(8)));
        else if (Q.indexOf("find ") === 0)
          G = find(asParentNode(J), normalizeSelector(Q.substr(5)));
        else if (Q === "next" || Q === "nextElementSibling")
          G = asElement(J).nextElementSibling;
        else if (Q.indexOf("next ") === 0)
          G = scanForwardQuery(J, normalizeSelector(Q.substr(5)), !!Z);
        else if (Q === "previous" || Q === "previousElementSibling")
          G = asElement(J).previousElementSibling;
        else if (Q.indexOf("previous ") === 0)
          G = scanBackwardsQuery(J, normalizeSelector(Q.substr(9)), !!Z);
        else if (Q === "document") G = document;
        else if (Q === "window") G = window;
        else if (Q === "body") G = document.body;
        else if (Q === "root") G = getRootNode(J, !!Z);
        else if (Q === "host") G = J.getRootNode().host;
        else X.push(Q);
        if (G) W.push(G);
      }
      if (X.length > 0) {
        let Q = X.join(","),
          G = asParentNode(getRootNode(J, !!Z));
        W.push(...toArray(G.querySelectorAll(Q)));
      }
      return W;
    }
    var scanForwardQuery = function (J, Y, Z) {
        let $ = asParentNode(getRootNode(J, Z)).querySelectorAll(Y);
        for (let W = 0; W < $.length; W++) {
          let X = $[W];
          if (X.compareDocumentPosition(J) === Node.DOCUMENT_POSITION_PRECEDING)
            return X;
        }
      },
      scanBackwardsQuery = function (J, Y, Z) {
        let $ = asParentNode(getRootNode(J, Z)).querySelectorAll(Y);
        for (let W = $.length - 1; W >= 0; W--) {
          let X = $[W];
          if (X.compareDocumentPosition(J) === Node.DOCUMENT_POSITION_FOLLOWING)
            return X;
        }
      };
    function querySelectorExt(J, Y) {
      if (typeof J !== "string") return querySelectorAllExt(J, Y)[0];
      else return querySelectorAllExt(getDocument().body, J)[0];
    }
    function resolveTarget(J, Y) {
      if (typeof J === "string") return find(asParentNode(Y) || document, J);
      else return J;
    }
    function processEventArgs(J, Y, Z, $) {
      if (isFunction(Y))
        return {
          target: getDocument().body,
          event: asString(J),
          listener: Y,
          options: Z,
        };
      else
        return {
          target: resolveTarget(J),
          event: asString(Y),
          listener: Z,
          options: $,
        };
    }
    function addEventListenerImpl(J, Y, Z, $) {
      return (
        ready(function () {
          let X = processEventArgs(J, Y, Z, $);
          X.target.addEventListener(X.event, X.listener, X.options);
        }),
        isFunction(Y) ? Y : Z
      );
    }
    function removeEventListenerImpl(J, Y, Z) {
      return (
        ready(function () {
          let $ = processEventArgs(J, Y, Z);
          $.target.removeEventListener($.event, $.listener);
        }),
        isFunction(Y) ? Y : Z
      );
    }
    let DUMMY_ELT = getDocument().createElement("output");
    function findAttributeTargets(J, Y) {
      let Z = getClosestAttributeValue(J, Y);
      if (Z)
        if (Z === "this") return [findThisElement(J, Y)];
        else {
          let $ = querySelectorAllExt(J, Z);
          if ($.length === 0)
            return (
              logError(
                'The selector "' + Z + '" on ' + Y + " returned no matches!",
              ),
              [DUMMY_ELT]
            );
          else return $;
        }
    }
    function findThisElement(J, Y) {
      return asElement(
        getClosestMatch(J, function (Z) {
          return getAttributeValue(asElement(Z), Y) != null;
        }),
      );
    }
    function getTarget(J) {
      let Y = getClosestAttributeValue(J, "hx-target");
      if (Y)
        if (Y === "this") return findThisElement(J, "hx-target");
        else return querySelectorExt(J, Y);
      else if (getInternalData(J).boosted) return getDocument().body;
      else return J;
    }
    function shouldSettleAttribute(J) {
      let Y = htmx.config.attributesToSettle;
      for (let Z = 0; Z < Y.length; Z++) if (J === Y[Z]) return !0;
      return !1;
    }
    function cloneAttributes(J, Y) {
      forEach(J.attributes, function (Z) {
        if (!Y.hasAttribute(Z.name) && shouldSettleAttribute(Z.name))
          J.removeAttribute(Z.name);
      }),
        forEach(Y.attributes, function (Z) {
          if (shouldSettleAttribute(Z.name)) J.setAttribute(Z.name, Z.value);
        });
    }
    function isInlineSwap(J, Y) {
      let Z = getExtensions(Y);
      for (let $ = 0; $ < Z.length; $++) {
        let W = Z[$];
        try {
          if (W.isInlineSwap(J)) return !0;
        } catch (X) {
          logError(X);
        }
      }
      return J === "outerHTML";
    }
    function oobSwap(J, Y, Z, $) {
      $ = $ || getDocument();
      let W = "#" + getRawAttribute(Y, "id"),
        X = "outerHTML";
      if (J === "true");
      else if (J.indexOf(":") > 0)
        (X = J.substring(0, J.indexOf(":"))),
          (W = J.substring(J.indexOf(":") + 1));
      else X = J;
      Y.removeAttribute("hx-swap-oob"), Y.removeAttribute("data-hx-swap-oob");
      let Q = querySelectorAllExt($, W, !1);
      if (Q)
        forEach(Q, function (G) {
          let B,
            U = Y.cloneNode(!0);
          if (
            ((B = getDocument().createDocumentFragment()),
            B.appendChild(U),
            !isInlineSwap(X, G))
          )
            B = asParentNode(U);
          let _ = { shouldSwap: !0, target: G, fragment: B };
          if (!triggerEvent(G, "htmx:oobBeforeSwap", _)) return;
          if (((G = _.target), _.shouldSwap))
            handlePreservedElements(B),
              swapWithStyle(X, G, G, B, Z),
              restorePreservedElements();
          forEach(Z.elts, function (z) {
            triggerEvent(z, "htmx:oobAfterSwap", _);
          });
        }),
          Y.parentNode.removeChild(Y);
      else
        Y.parentNode.removeChild(Y),
          triggerErrorEvent(getDocument().body, "htmx:oobErrorNoTarget", {
            content: Y,
          });
      return J;
    }
    function restorePreservedElements() {
      let J = find("#--htmx-preserve-pantry--");
      if (J) {
        for (let Y of [...J.children]) {
          let Z = find("#" + Y.id);
          Z.parentNode.moveBefore(Y, Z), Z.remove();
        }
        J.remove();
      }
    }
    function handlePreservedElements(J) {
      forEach(findAll(J, "[hx-preserve], [data-hx-preserve]"), function (Y) {
        let Z = getAttributeValue(Y, "id"),
          $ = getDocument().getElementById(Z);
        if ($ != null)
          if (Y.moveBefore) {
            let W = find("#--htmx-preserve-pantry--");
            if (W == null)
              getDocument().body.insertAdjacentHTML(
                "afterend",
                "<div id='--htmx-preserve-pantry--'></div>",
              ),
                (W = find("#--htmx-preserve-pantry--"));
            W.moveBefore($, null);
          } else Y.parentNode.replaceChild($, Y);
      });
    }
    function handleAttributes(J, Y, Z) {
      forEach(Y.querySelectorAll("[id]"), function ($) {
        let W = getRawAttribute($, "id");
        if (W && W.length > 0) {
          let X = W.replace("'", "\\'"),
            Q = $.tagName.replace(":", "\\:"),
            G = asParentNode(J),
            B = G && G.querySelector(Q + "[id='" + X + "']");
          if (B && B !== G) {
            let U = $.cloneNode();
            cloneAttributes($, B),
              Z.tasks.push(function () {
                cloneAttributes($, U);
              });
          }
        }
      });
    }
    function makeAjaxLoadTask(J) {
      return function () {
        removeClassFromElement(J, htmx.config.addedClass),
          processNode(asElement(J)),
          processFocus(asParentNode(J)),
          triggerEvent(J, "htmx:load");
      };
    }
    function processFocus(J) {
      let Z = asHtmlElement(
        matches(J, "[autofocus]") ? J : J.querySelector("[autofocus]"),
      );
      if (Z != null) Z.focus();
    }
    function insertNodesBefore(J, Y, Z, $) {
      handleAttributes(J, Z, $);
      while (Z.childNodes.length > 0) {
        let W = Z.firstChild;
        if (
          (addClassToElement(asElement(W), htmx.config.addedClass),
          J.insertBefore(W, Y),
          W.nodeType !== Node.TEXT_NODE && W.nodeType !== Node.COMMENT_NODE)
        )
          $.tasks.push(makeAjaxLoadTask(W));
      }
    }
    function stringHash(J, Y) {
      let Z = 0;
      while (Z < J.length) Y = ((Y << 5) - Y + J.charCodeAt(Z++)) | 0;
      return Y;
    }
    function attributeHash(J) {
      let Y = 0;
      if (J.attributes)
        for (let Z = 0; Z < J.attributes.length; Z++) {
          let $ = J.attributes[Z];
          if ($.value)
            (Y = stringHash($.name, Y)), (Y = stringHash($.value, Y));
        }
      return Y;
    }
    function deInitOnHandlers(J) {
      let Y = getInternalData(J);
      if (Y.onHandlers) {
        for (let Z = 0; Z < Y.onHandlers.length; Z++) {
          let $ = Y.onHandlers[Z];
          removeEventListenerImpl(J, $.event, $.listener);
        }
        delete Y.onHandlers;
      }
    }
    function deInitNode(J) {
      let Y = getInternalData(J);
      if (Y.timeout) clearTimeout(Y.timeout);
      if (Y.listenerInfos)
        forEach(Y.listenerInfos, function (Z) {
          if (Z.on) removeEventListenerImpl(Z.on, Z.trigger, Z.listener);
        });
      deInitOnHandlers(J),
        forEach(Object.keys(Y), function (Z) {
          if (Z !== "firstInitCompleted") delete Y[Z];
        });
    }
    function cleanUpElement(J) {
      if (
        (triggerEvent(J, "htmx:beforeCleanupElement"),
        deInitNode(J),
        J.children)
      )
        forEach(J.children, function (Y) {
          cleanUpElement(Y);
        });
    }
    function swapOuterHTML(J, Y, Z) {
      if (J instanceof Element && J.tagName === "BODY")
        return swapInnerHTML(J, Y, Z);
      let $,
        W = J.previousSibling,
        X = parentElt(J);
      if (!X) return;
      if ((insertNodesBefore(X, J, Y, Z), W == null)) $ = X.firstChild;
      else $ = W.nextSibling;
      Z.elts = Z.elts.filter(function (Q) {
        return Q !== J;
      });
      while ($ && $ !== J) {
        if ($ instanceof Element) Z.elts.push($);
        $ = $.nextSibling;
      }
      if ((cleanUpElement(J), J instanceof Element)) J.remove();
      else J.parentNode.removeChild(J);
    }
    function swapAfterBegin(J, Y, Z) {
      return insertNodesBefore(J, J.firstChild, Y, Z);
    }
    function swapBeforeBegin(J, Y, Z) {
      return insertNodesBefore(parentElt(J), J, Y, Z);
    }
    function swapBeforeEnd(J, Y, Z) {
      return insertNodesBefore(J, null, Y, Z);
    }
    function swapAfterEnd(J, Y, Z) {
      return insertNodesBefore(parentElt(J), J.nextSibling, Y, Z);
    }
    function swapDelete(J) {
      cleanUpElement(J);
      let Y = parentElt(J);
      if (Y) return Y.removeChild(J);
    }
    function swapInnerHTML(J, Y, Z) {
      let $ = J.firstChild;
      if ((insertNodesBefore(J, $, Y, Z), $)) {
        while ($.nextSibling)
          cleanUpElement($.nextSibling), J.removeChild($.nextSibling);
        cleanUpElement($), J.removeChild($);
      }
    }
    function swapWithStyle(J, Y, Z, $, W) {
      switch (J) {
        case "none":
          return;
        case "outerHTML":
          swapOuterHTML(Z, $, W);
          return;
        case "afterbegin":
          swapAfterBegin(Z, $, W);
          return;
        case "beforebegin":
          swapBeforeBegin(Z, $, W);
          return;
        case "beforeend":
          swapBeforeEnd(Z, $, W);
          return;
        case "afterend":
          swapAfterEnd(Z, $, W);
          return;
        case "delete":
          swapDelete(Z);
          return;
        default:
          var X = getExtensions(Y);
          for (let Q = 0; Q < X.length; Q++) {
            let G = X[Q];
            try {
              let B = G.handleSwap(J, Z, $, W);
              if (B) {
                if (Array.isArray(B))
                  for (let U = 0; U < B.length; U++) {
                    let _ = B[U];
                    if (
                      _.nodeType !== Node.TEXT_NODE &&
                      _.nodeType !== Node.COMMENT_NODE
                    )
                      W.tasks.push(makeAjaxLoadTask(_));
                  }
                return;
              }
            } catch (B) {
              logError(B);
            }
          }
          if (J === "innerHTML") swapInnerHTML(Z, $, W);
          else swapWithStyle(htmx.config.defaultSwapStyle, Y, Z, $, W);
      }
    }
    function findAndSwapOobElements(J, Y, Z) {
      var $ = findAll(J, "[hx-swap-oob], [data-hx-swap-oob]");
      return (
        forEach($, function (W) {
          if (htmx.config.allowNestedOobSwaps || W.parentElement === null) {
            let X = getAttributeValue(W, "hx-swap-oob");
            if (X != null) oobSwap(X, W, Y, Z);
          } else
            W.removeAttribute("hx-swap-oob"),
              W.removeAttribute("data-hx-swap-oob");
        }),
        $.length > 0
      );
    }
    function swap(J, Y, Z, $) {
      if (!$) $ = {};
      J = resolveTarget(J);
      let W = $.contextElement
          ? getRootNode($.contextElement, !1)
          : getDocument(),
        X = document.activeElement,
        Q = {};
      try {
        Q = {
          elt: X,
          start: X ? X.selectionStart : null,
          end: X ? X.selectionEnd : null,
        };
      } catch (U) {}
      let G = makeSettleInfo(J);
      if (Z.swapStyle === "textContent") J.textContent = Y;
      else {
        let U = makeFragment(Y);
        if (((G.title = U.title), $.selectOOB)) {
          let _ = $.selectOOB.split(",");
          for (let z = 0; z < _.length; z++) {
            let K = _[z].split(":", 2),
              q = K[0].trim();
            if (q.indexOf("#") === 0) q = q.substring(1);
            let j = K[1] || "true",
              T = U.querySelector("#" + q);
            if (T) oobSwap(j, T, G, W);
          }
        }
        if (
          (findAndSwapOobElements(U, G, W),
          forEach(findAll(U, "template"), function (_) {
            if (_.content && findAndSwapOobElements(_.content, G, W))
              _.remove();
          }),
          $.select)
        ) {
          let _ = getDocument().createDocumentFragment();
          forEach(U.querySelectorAll($.select), function (z) {
            _.appendChild(z);
          }),
            (U = _);
        }
        handlePreservedElements(U),
          swapWithStyle(Z.swapStyle, $.contextElement, J, U, G),
          restorePreservedElements();
      }
      if (Q.elt && !bodyContains(Q.elt) && getRawAttribute(Q.elt, "id")) {
        let U = document.getElementById(getRawAttribute(Q.elt, "id")),
          _ = {
            preventScroll:
              Z.focusScroll !== void 0
                ? !Z.focusScroll
                : !htmx.config.defaultFocusScroll,
          };
        if (U) {
          if (Q.start && U.setSelectionRange)
            try {
              U.setSelectionRange(Q.start, Q.end);
            } catch (z) {}
          U.focus(_);
        }
      }
      if (
        (J.classList.remove(htmx.config.swappingClass),
        forEach(G.elts, function (U) {
          if (U.classList) U.classList.add(htmx.config.settlingClass);
          triggerEvent(U, "htmx:afterSwap", $.eventInfo);
        }),
        $.afterSwapCallback)
      )
        $.afterSwapCallback();
      if (!Z.ignoreTitle) handleTitle(G.title);
      let B = function () {
        if (
          (forEach(G.tasks, function (U) {
            U.call();
          }),
          forEach(G.elts, function (U) {
            if (U.classList) U.classList.remove(htmx.config.settlingClass);
            triggerEvent(U, "htmx:afterSettle", $.eventInfo);
          }),
          $.anchor)
        ) {
          let U = asElement(resolveTarget("#" + $.anchor));
          if (U) U.scrollIntoView({ block: "start", behavior: "auto" });
        }
        if ((updateScrollState(G.elts, Z), $.afterSettleCallback))
          $.afterSettleCallback();
      };
      if (Z.settleDelay > 0) getWindow().setTimeout(B, Z.settleDelay);
      else B();
    }
    function handleTriggerHeader(J, Y, Z) {
      let $ = J.getResponseHeader(Y);
      if ($.indexOf("{") === 0) {
        let W = parseJSON($);
        for (let X in W)
          if (W.hasOwnProperty(X)) {
            let Q = W[X];
            if (isRawObject(Q)) Z = Q.target !== void 0 ? Q.target : Z;
            else Q = { value: Q };
            triggerEvent(Z, X, Q);
          }
      } else {
        let W = $.split(",");
        for (let X = 0; X < W.length; X++) triggerEvent(Z, W[X].trim(), []);
      }
    }
    let WHITESPACE = /\s/,
      WHITESPACE_OR_COMMA = /[\s,]/,
      SYMBOL_START = /[_$a-zA-Z]/,
      SYMBOL_CONT = /[_$a-zA-Z0-9]/,
      STRINGISH_START = ['"', "'", "/"],
      NOT_WHITESPACE = /[^\s]/,
      COMBINED_SELECTOR_START = /[{(]/,
      COMBINED_SELECTOR_END = /[})]/;
    function tokenizeString(J) {
      let Y = [],
        Z = 0;
      while (Z < J.length) {
        if (SYMBOL_START.exec(J.charAt(Z))) {
          var $ = Z;
          while (SYMBOL_CONT.exec(J.charAt(Z + 1))) Z++;
          Y.push(J.substring($, Z + 1));
        } else if (STRINGISH_START.indexOf(J.charAt(Z)) !== -1) {
          let W = J.charAt(Z);
          var $ = Z;
          Z++;
          while (Z < J.length && J.charAt(Z) !== W) {
            if (J.charAt(Z) === "\\") Z++;
            Z++;
          }
          Y.push(J.substring($, Z + 1));
        } else {
          let W = J.charAt(Z);
          Y.push(W);
        }
        Z++;
      }
      return Y;
    }
    function isPossibleRelativeReference(J, Y, Z) {
      return (
        SYMBOL_START.exec(J.charAt(0)) &&
        J !== "true" &&
        J !== "false" &&
        J !== "this" &&
        J !== Z &&
        Y !== "."
      );
    }
    function maybeGenerateConditional(J, Y, Z) {
      if (Y[0] === "[") {
        Y.shift();
        let $ = 1,
          W = " return (function(" + Z + "){ return (",
          X = null;
        while (Y.length > 0) {
          let Q = Y[0];
          if (Q === "]") {
            if (($--, $ === 0)) {
              if (X === null) W = W + "true";
              Y.shift(), (W += ")})");
              try {
                let G = maybeEval(
                  J,
                  function () {
                    return Function(W)();
                  },
                  function () {
                    return !0;
                  },
                );
                return (G.source = W), G;
              } catch (G) {
                return (
                  triggerErrorEvent(getDocument().body, "htmx:syntax:error", {
                    error: G,
                    source: W,
                  }),
                  null
                );
              }
            }
          } else if (Q === "[") $++;
          if (isPossibleRelativeReference(Q, X, Z))
            W +=
              "((" +
              Z +
              "." +
              Q +
              ") ? (" +
              Z +
              "." +
              Q +
              ") : (window." +
              Q +
              "))";
          else W = W + Q;
          X = Y.shift();
        }
      }
    }
    function consumeUntil(J, Y) {
      let Z = "";
      while (J.length > 0 && !Y.test(J[0])) Z += J.shift();
      return Z;
    }
    function consumeCSSSelector(J) {
      let Y;
      if (J.length > 0 && COMBINED_SELECTOR_START.test(J[0]))
        J.shift(),
          (Y = consumeUntil(J, COMBINED_SELECTOR_END).trim()),
          J.shift();
      else Y = consumeUntil(J, WHITESPACE_OR_COMMA);
      return Y;
    }
    let INPUT_SELECTOR = "input, textarea, select";
    function parseAndCacheTrigger(J, Y, Z) {
      let $ = [],
        W = tokenizeString(Y);
      do {
        consumeUntil(W, NOT_WHITESPACE);
        let G = W.length,
          B = consumeUntil(W, /[,\[\s]/);
        if (B !== "")
          if (B === "every") {
            let U = { trigger: "every" };
            consumeUntil(W, NOT_WHITESPACE),
              (U.pollInterval = parseInterval(consumeUntil(W, /[,\[\s]/))),
              consumeUntil(W, NOT_WHITESPACE);
            var X = maybeGenerateConditional(J, W, "event");
            if (X) U.eventFilter = X;
            $.push(U);
          } else {
            let U = { trigger: B };
            var X = maybeGenerateConditional(J, W, "event");
            if (X) U.eventFilter = X;
            consumeUntil(W, NOT_WHITESPACE);
            while (W.length > 0 && W[0] !== ",") {
              let z = W.shift();
              if (z === "changed") U.changed = !0;
              else if (z === "once") U.once = !0;
              else if (z === "consume") U.consume = !0;
              else if (z === "delay" && W[0] === ":")
                W.shift(),
                  (U.delay = parseInterval(
                    consumeUntil(W, WHITESPACE_OR_COMMA),
                  ));
              else if (z === "from" && W[0] === ":") {
                if ((W.shift(), COMBINED_SELECTOR_START.test(W[0])))
                  var Q = consumeCSSSelector(W);
                else {
                  var Q = consumeUntil(W, WHITESPACE_OR_COMMA);
                  if (
                    Q === "closest" ||
                    Q === "find" ||
                    Q === "next" ||
                    Q === "previous"
                  ) {
                    W.shift();
                    let q = consumeCSSSelector(W);
                    if (q.length > 0) Q += " " + q;
                  }
                }
                U.from = Q;
              } else if (z === "target" && W[0] === ":")
                W.shift(), (U.target = consumeCSSSelector(W));
              else if (z === "throttle" && W[0] === ":")
                W.shift(),
                  (U.throttle = parseInterval(
                    consumeUntil(W, WHITESPACE_OR_COMMA),
                  ));
              else if (z === "queue" && W[0] === ":")
                W.shift(), (U.queue = consumeUntil(W, WHITESPACE_OR_COMMA));
              else if (z === "root" && W[0] === ":")
                W.shift(), (U[z] = consumeCSSSelector(W));
              else if (z === "threshold" && W[0] === ":")
                W.shift(), (U[z] = consumeUntil(W, WHITESPACE_OR_COMMA));
              else
                triggerErrorEvent(J, "htmx:syntax:error", { token: W.shift() });
              consumeUntil(W, NOT_WHITESPACE);
            }
            $.push(U);
          }
        if (W.length === G)
          triggerErrorEvent(J, "htmx:syntax:error", { token: W.shift() });
        consumeUntil(W, NOT_WHITESPACE);
      } while (W[0] === "," && W.shift());
      if (Z) Z[Y] = $;
      return $;
    }
    function getTriggerSpecs(J) {
      let Y = getAttributeValue(J, "hx-trigger"),
        Z = [];
      if (Y) {
        let $ = htmx.config.triggerSpecsCache;
        Z = ($ && $[Y]) || parseAndCacheTrigger(J, Y, $);
      }
      if (Z.length > 0) return Z;
      else if (matches(J, "form")) return [{ trigger: "submit" }];
      else if (matches(J, 'input[type="button"], input[type="submit"]'))
        return [{ trigger: "click" }];
      else if (matches(J, INPUT_SELECTOR)) return [{ trigger: "change" }];
      else return [{ trigger: "click" }];
    }
    function cancelPolling(J) {
      getInternalData(J).cancelled = !0;
    }
    function processPolling(J, Y, Z) {
      let $ = getInternalData(J);
      $.timeout = getWindow().setTimeout(function () {
        if (bodyContains(J) && $.cancelled !== !0) {
          if (
            !maybeFilterEvent(
              Z,
              J,
              makeEvent("hx:poll:trigger", { triggerSpec: Z, target: J }),
            )
          )
            Y(J);
          processPolling(J, Y, Z);
        }
      }, Z.pollInterval);
    }
    function isLocalLink(J) {
      return (
        location.hostname === J.hostname &&
        getRawAttribute(J, "href") &&
        getRawAttribute(J, "href").indexOf("#") !== 0
      );
    }
    function eltIsDisabled(J) {
      return closest(J, htmx.config.disableSelector);
    }
    function boostElement(J, Y, Z) {
      if (
        (J instanceof HTMLAnchorElement &&
          isLocalLink(J) &&
          (J.target === "" || J.target === "_self")) ||
        (J.tagName === "FORM" &&
          String(getRawAttribute(J, "method")).toLowerCase() !== "dialog")
      ) {
        Y.boosted = !0;
        let $, W;
        if (J.tagName === "A") ($ = "get"), (W = getRawAttribute(J, "href"));
        else {
          let X = getRawAttribute(J, "method");
          if (
            (($ = X ? X.toLowerCase() : "get"),
            (W = getRawAttribute(J, "action")),
            W == null || W === "")
          )
            W = getDocument().location.href;
          if ($ === "get" && W.includes("?")) W = W.replace(/\?[^#]+/, "");
        }
        Z.forEach(function (X) {
          addEventListener(
            J,
            function (Q, G) {
              let B = asElement(Q);
              if (eltIsDisabled(B)) {
                cleanUpElement(B);
                return;
              }
              issueAjaxRequest($, W, B, G);
            },
            Y,
            X,
            !0,
          );
        });
      }
    }
    function shouldCancel(J, Y) {
      let Z = asElement(Y);
      if (!Z) return !1;
      if (J.type === "submit" || J.type === "click") {
        if (Z.tagName === "FORM") return !0;
        if (
          matches(Z, 'input[type="submit"], button') &&
          (matches(Z, "[form]") || closest(Z, "form") !== null)
        )
          return !0;
        if (
          Z instanceof HTMLAnchorElement &&
          Z.href &&
          (Z.getAttribute("href") === "#" ||
            Z.getAttribute("href").indexOf("#") !== 0)
        )
          return !0;
      }
      return !1;
    }
    function ignoreBoostedAnchorCtrlClick(J, Y) {
      return (
        getInternalData(J).boosted &&
        J instanceof HTMLAnchorElement &&
        Y.type === "click" &&
        (Y.ctrlKey || Y.metaKey)
      );
    }
    function maybeFilterEvent(J, Y, Z) {
      let $ = J.eventFilter;
      if ($)
        try {
          return $.call(Y, Z) !== !0;
        } catch (W) {
          let X = $.source;
          return (
            triggerErrorEvent(getDocument().body, "htmx:eventFilter:error", {
              error: W,
              source: X,
            }),
            !0
          );
        }
      return !1;
    }
    function addEventListener(J, Y, Z, $, W) {
      let X = getInternalData(J),
        Q;
      if ($.from) Q = querySelectorAllExt(J, $.from);
      else Q = [J];
      if ($.changed) {
        if (!("lastValue" in X)) X.lastValue = new WeakMap();
        Q.forEach(function (G) {
          if (!X.lastValue.has($)) X.lastValue.set($, new WeakMap());
          X.lastValue.get($).set(G, G.value);
        });
      }
      forEach(Q, function (G) {
        let B = function (U) {
          if (!bodyContains(J)) {
            G.removeEventListener($.trigger, B);
            return;
          }
          if (ignoreBoostedAnchorCtrlClick(J, U)) return;
          if (W || shouldCancel(U, J)) U.preventDefault();
          if (maybeFilterEvent($, J, U)) return;
          let _ = getInternalData(U);
          if (((_.triggerSpec = $), _.handledFor == null)) _.handledFor = [];
          if (_.handledFor.indexOf(J) < 0) {
            if ((_.handledFor.push(J), $.consume)) U.stopPropagation();
            if ($.target && U.target) {
              if (!matches(asElement(U.target), $.target)) return;
            }
            if ($.once)
              if (X.triggeredOnce) return;
              else X.triggeredOnce = !0;
            if ($.changed) {
              let z = event.target,
                K = z.value,
                q = X.lastValue.get($);
              if (q.has(z) && q.get(z) === K) return;
              q.set(z, K);
            }
            if (X.delayed) clearTimeout(X.delayed);
            if (X.throttle) return;
            if ($.throttle > 0) {
              if (!X.throttle)
                triggerEvent(J, "htmx:trigger"),
                  Y(J, U),
                  (X.throttle = getWindow().setTimeout(function () {
                    X.throttle = null;
                  }, $.throttle));
            } else if ($.delay > 0)
              X.delayed = getWindow().setTimeout(function () {
                triggerEvent(J, "htmx:trigger"), Y(J, U);
              }, $.delay);
            else triggerEvent(J, "htmx:trigger"), Y(J, U);
          }
        };
        if (Z.listenerInfos == null) Z.listenerInfos = [];
        Z.listenerInfos.push({ trigger: $.trigger, listener: B, on: G }),
          G.addEventListener($.trigger, B);
      });
    }
    let windowIsScrolling = !1,
      scrollHandler = null;
    function initScrollHandler() {
      if (!scrollHandler)
        (scrollHandler = function () {
          windowIsScrolling = !0;
        }),
          window.addEventListener("scroll", scrollHandler),
          window.addEventListener("resize", scrollHandler),
          setInterval(function () {
            if (windowIsScrolling)
              (windowIsScrolling = !1),
                forEach(
                  getDocument().querySelectorAll(
                    "[hx-trigger*='revealed'],[data-hx-trigger*='revealed']",
                  ),
                  function (J) {
                    maybeReveal(J);
                  },
                );
          }, 200);
    }
    function maybeReveal(J) {
      if (!hasAttribute(J, "data-hx-revealed") && isScrolledIntoView(J))
        if (
          (J.setAttribute("data-hx-revealed", "true"),
          getInternalData(J).initHash)
        )
          triggerEvent(J, "revealed");
        else
          J.addEventListener(
            "htmx:afterProcessNode",
            function () {
              triggerEvent(J, "revealed");
            },
            { once: !0 },
          );
    }
    function loadImmediately(J, Y, Z, $) {
      let W = function () {
        if (!Z.loaded) (Z.loaded = !0), triggerEvent(J, "htmx:trigger"), Y(J);
      };
      if ($ > 0) getWindow().setTimeout(W, $);
      else W();
    }
    function processVerbs(J, Y, Z) {
      let $ = !1;
      return (
        forEach(VERBS, function (W) {
          if (hasAttribute(J, "hx-" + W)) {
            let X = getAttributeValue(J, "hx-" + W);
            ($ = !0),
              (Y.path = X),
              (Y.verb = W),
              Z.forEach(function (Q) {
                addTriggerHandler(J, Q, Y, function (G, B) {
                  let U = asElement(G);
                  if (closest(U, htmx.config.disableSelector)) {
                    cleanUpElement(U);
                    return;
                  }
                  issueAjaxRequest(W, X, U, B);
                });
              });
          }
        }),
        $
      );
    }
    function addTriggerHandler(J, Y, Z, $) {
      if (Y.trigger === "revealed")
        initScrollHandler(),
          addEventListener(J, $, Z, Y),
          maybeReveal(asElement(J));
      else if (Y.trigger === "intersect") {
        let W = {};
        if (Y.root) W.root = querySelectorExt(J, Y.root);
        if (Y.threshold) W.threshold = parseFloat(Y.threshold);
        new IntersectionObserver(function (Q) {
          for (let G = 0; G < Q.length; G++)
            if (Q[G].isIntersecting) {
              triggerEvent(J, "intersect");
              break;
            }
        }, W).observe(asElement(J)),
          addEventListener(asElement(J), $, Z, Y);
      } else if (!Z.firstInitCompleted && Y.trigger === "load") {
        if (!maybeFilterEvent(Y, J, makeEvent("load", { elt: J })))
          loadImmediately(asElement(J), $, Z, Y.delay);
      } else if (Y.pollInterval > 0)
        (Z.polling = !0), processPolling(asElement(J), $, Y);
      else addEventListener(J, $, Z, Y);
    }
    function shouldProcessHxOn(J) {
      let Y = asElement(J);
      if (!Y) return !1;
      let Z = Y.attributes;
      for (let $ = 0; $ < Z.length; $++) {
        let W = Z[$].name;
        if (
          startsWith(W, "hx-on:") ||
          startsWith(W, "data-hx-on:") ||
          startsWith(W, "hx-on-") ||
          startsWith(W, "data-hx-on-")
        )
          return !0;
      }
      return !1;
    }
    let HX_ON_QUERY = new XPathEvaluator().createExpression(
      './/*[@*[ starts-with(name(), "hx-on:") or starts-with(name(), "data-hx-on:") or starts-with(name(), "hx-on-") or starts-with(name(), "data-hx-on-") ]]',
    );
    function processHXOnRoot(J, Y) {
      if (shouldProcessHxOn(J)) Y.push(asElement(J));
      let Z = HX_ON_QUERY.evaluate(J),
        $ = null;
      while (($ = Z.iterateNext())) Y.push(asElement($));
    }
    function findHxOnWildcardElements(J) {
      let Y = [];
      if (J instanceof DocumentFragment)
        for (let Z of J.childNodes) processHXOnRoot(Z, Y);
      else processHXOnRoot(J, Y);
      return Y;
    }
    function findElementsToProcess(J) {
      if (J.querySelectorAll) {
        let $ = [];
        for (let X in extensions) {
          let Q = extensions[X];
          if (Q.getSelectors) {
            var Y = Q.getSelectors();
            if (Y) $.push(Y);
          }
        }
        return J.querySelectorAll(
          VERB_SELECTOR +
            ", [hx-boost] a, [data-hx-boost] a, a[hx-boost], a[data-hx-boost], form, [type='submit'], [hx-ext], [data-hx-ext], [hx-trigger], [data-hx-trigger]" +
            $.flat()
              .map((X) => ", " + X)
              .join(""),
        );
      } else return [];
    }
    function maybeSetLastButtonClicked(J) {
      let Y = closest(asElement(J.target), "button, input[type='submit']"),
        Z = getRelatedFormData(J);
      if (Z) Z.lastButtonClicked = Y;
    }
    function maybeUnsetLastButtonClicked(J) {
      let Y = getRelatedFormData(J);
      if (Y) Y.lastButtonClicked = null;
    }
    function getRelatedFormData(J) {
      let Y = closest(asElement(J.target), "button, input[type='submit']");
      if (!Y) return;
      let Z =
        resolveTarget("#" + getRawAttribute(Y, "form"), Y.getRootNode()) ||
        closest(Y, "form");
      if (!Z) return;
      return getInternalData(Z);
    }
    function initButtonTracking(J) {
      J.addEventListener("click", maybeSetLastButtonClicked),
        J.addEventListener("focusin", maybeSetLastButtonClicked),
        J.addEventListener("focusout", maybeUnsetLastButtonClicked);
    }
    function addHxOnEventHandler(J, Y, Z) {
      let $ = getInternalData(J);
      if (!Array.isArray($.onHandlers)) $.onHandlers = [];
      let W,
        X = function (Q) {
          maybeEval(J, function () {
            if (eltIsDisabled(J)) return;
            if (!W) W = new Function("event", Z);
            W.call(J, Q);
          });
        };
      J.addEventListener(Y, X), $.onHandlers.push({ event: Y, listener: X });
    }
    function processHxOnWildcard(J) {
      deInitOnHandlers(J);
      for (let Y = 0; Y < J.attributes.length; Y++) {
        let Z = J.attributes[Y].name,
          $ = J.attributes[Y].value;
        if (startsWith(Z, "hx-on") || startsWith(Z, "data-hx-on")) {
          let W = Z.indexOf("-on") + 3,
            X = Z.slice(W, W + 1);
          if (X === "-" || X === ":") {
            let Q = Z.slice(W + 1);
            if (startsWith(Q, ":")) Q = "htmx" + Q;
            else if (startsWith(Q, "-")) Q = "htmx:" + Q.slice(1);
            else if (startsWith(Q, "htmx-")) Q = "htmx:" + Q.slice(5);
            addHxOnEventHandler(J, Q, $);
          }
        }
      }
    }
    function initNode(J) {
      if (closest(J, htmx.config.disableSelector)) {
        cleanUpElement(J);
        return;
      }
      let Y = getInternalData(J),
        Z = attributeHash(J);
      if (Y.initHash !== Z) {
        deInitNode(J),
          (Y.initHash = Z),
          triggerEvent(J, "htmx:beforeProcessNode");
        let $ = getTriggerSpecs(J);
        if (!processVerbs(J, Y, $)) {
          if (getClosestAttributeValue(J, "hx-boost") === "true")
            boostElement(J, Y, $);
          else if (hasAttribute(J, "hx-trigger"))
            $.forEach(function (X) {
              addTriggerHandler(J, X, Y, function () {});
            });
        }
        if (
          J.tagName === "FORM" ||
          (getRawAttribute(J, "type") === "submit" && hasAttribute(J, "form"))
        )
          initButtonTracking(J);
        (Y.firstInitCompleted = !0), triggerEvent(J, "htmx:afterProcessNode");
      }
    }
    function processNode(J) {
      if (((J = resolveTarget(J)), closest(J, htmx.config.disableSelector))) {
        cleanUpElement(J);
        return;
      }
      initNode(J),
        forEach(findElementsToProcess(J), function (Y) {
          initNode(Y);
        }),
        forEach(findHxOnWildcardElements(J), processHxOnWildcard);
    }
    function kebabEventName(J) {
      return J.replace(/([a-z0-9])([A-Z])/g, "$1-$2").toLowerCase();
    }
    function makeEvent(J, Y) {
      let Z;
      if (window.CustomEvent && typeof window.CustomEvent === "function")
        Z = new CustomEvent(J, {
          bubbles: !0,
          cancelable: !0,
          composed: !0,
          detail: Y,
        });
      else
        (Z = getDocument().createEvent("CustomEvent")),
          Z.initCustomEvent(J, !0, !0, Y);
      return Z;
    }
    function triggerErrorEvent(J, Y, Z) {
      triggerEvent(J, Y, mergeObjects({ error: Y }, Z));
    }
    function ignoreEventForLogging(J) {
      return J === "htmx:afterProcessNode";
    }
    function withExtensions(J, Y) {
      forEach(getExtensions(J), function (Z) {
        try {
          Y(Z);
        } catch ($) {
          logError($);
        }
      });
    }
    function logError(J) {
      if (console.error) console.error(J);
      else if (console.log) console.log("ERROR: ", J);
    }
    function triggerEvent(J, Y, Z) {
      if (((J = resolveTarget(J)), Z == null)) Z = {};
      Z.elt = J;
      let $ = makeEvent(Y, Z);
      if (htmx.logger && !ignoreEventForLogging(Y)) htmx.logger(J, Y, Z);
      if (Z.error)
        logError(Z.error), triggerEvent(J, "htmx:error", { errorInfo: Z });
      let W = J.dispatchEvent($),
        X = kebabEventName(Y);
      if (W && X !== Y) {
        let Q = makeEvent(X, $.detail);
        W = W && J.dispatchEvent(Q);
      }
      return (
        withExtensions(asElement(J), function (Q) {
          W = W && Q.onEvent(Y, $) !== !1 && !$.defaultPrevented;
        }),
        W
      );
    }
    let currentPathForHistory = location.pathname + location.search;
    function getHistoryElement() {
      return (
        getDocument().querySelector("[hx-history-elt],[data-hx-history-elt]") ||
        getDocument().body
      );
    }
    function saveToHistoryCache(J, Y) {
      if (!canAccessLocalStorage()) return;
      let Z = cleanInnerHtmlForHistory(Y),
        $ = getDocument().title,
        W = window.scrollY;
      if (htmx.config.historyCacheSize <= 0) {
        localStorage.removeItem("htmx-history-cache");
        return;
      }
      J = normalizePath(J);
      let X = parseJSON(localStorage.getItem("htmx-history-cache")) || [];
      for (let G = 0; G < X.length; G++)
        if (X[G].url === J) {
          X.splice(G, 1);
          break;
        }
      let Q = { url: J, content: Z, title: $, scroll: W };
      triggerEvent(getDocument().body, "htmx:historyItemCreated", {
        item: Q,
        cache: X,
      }),
        X.push(Q);
      while (X.length > htmx.config.historyCacheSize) X.shift();
      while (X.length > 0)
        try {
          localStorage.setItem("htmx-history-cache", JSON.stringify(X));
          break;
        } catch (G) {
          triggerErrorEvent(getDocument().body, "htmx:historyCacheError", {
            cause: G,
            cache: X,
          }),
            X.shift();
        }
    }
    function getCachedHistory(J) {
      if (!canAccessLocalStorage()) return null;
      J = normalizePath(J);
      let Y = parseJSON(localStorage.getItem("htmx-history-cache")) || [];
      for (let Z = 0; Z < Y.length; Z++) if (Y[Z].url === J) return Y[Z];
      return null;
    }
    function cleanInnerHtmlForHistory(J) {
      let Y = htmx.config.requestClass,
        Z = J.cloneNode(!0);
      return (
        forEach(findAll(Z, "." + Y), function ($) {
          removeClassFromElement($, Y);
        }),
        forEach(findAll(Z, "[data-disabled-by-htmx]"), function ($) {
          $.removeAttribute("disabled");
        }),
        Z.innerHTML
      );
    }
    function saveCurrentPageToHistory() {
      let J = getHistoryElement(),
        Y = currentPathForHistory || location.pathname + location.search,
        Z;
      try {
        Z = getDocument().querySelector(
          '[hx-history="false" i],[data-hx-history="false" i]',
        );
      } catch ($) {
        Z = getDocument().querySelector(
          '[hx-history="false"],[data-hx-history="false"]',
        );
      }
      if (!Z)
        triggerEvent(getDocument().body, "htmx:beforeHistorySave", {
          path: Y,
          historyElt: J,
        }),
          saveToHistoryCache(Y, J);
      if (htmx.config.historyEnabled)
        history.replaceState(
          { htmx: !0 },
          getDocument().title,
          window.location.href,
        );
    }
    function pushUrlIntoHistory(J) {
      if (htmx.config.getCacheBusterParam) {
        if (
          ((J = J.replace(/org\.htmx\.cache-buster=[^&]*&?/, "")),
          endsWith(J, "&") || endsWith(J, "?"))
        )
          J = J.slice(0, -1);
      }
      if (htmx.config.historyEnabled) history.pushState({ htmx: !0 }, "", J);
      currentPathForHistory = J;
    }
    function replaceUrlInHistory(J) {
      if (htmx.config.historyEnabled) history.replaceState({ htmx: !0 }, "", J);
      currentPathForHistory = J;
    }
    function settleImmediately(J) {
      forEach(J, function (Y) {
        Y.call(void 0);
      });
    }
    function loadHistoryFromServer(J) {
      let Y = new XMLHttpRequest(),
        Z = { path: J, xhr: Y };
      triggerEvent(getDocument().body, "htmx:historyCacheMiss", Z),
        Y.open("GET", J, !0),
        Y.setRequestHeader("HX-Request", "true"),
        Y.setRequestHeader("HX-History-Restore-Request", "true"),
        Y.setRequestHeader("HX-Current-URL", getDocument().location.href),
        (Y.onload = function () {
          if (this.status >= 200 && this.status < 400) {
            triggerEvent(getDocument().body, "htmx:historyCacheMissLoad", Z);
            let $ = makeFragment(this.response),
              W =
                $.querySelector("[hx-history-elt],[data-hx-history-elt]") || $,
              X = getHistoryElement(),
              Q = makeSettleInfo(X);
            handleTitle($.title),
              handlePreservedElements($),
              swapInnerHTML(X, W, Q),
              restorePreservedElements(),
              settleImmediately(Q.tasks),
              (currentPathForHistory = J),
              triggerEvent(getDocument().body, "htmx:historyRestore", {
                path: J,
                cacheMiss: !0,
                serverResponse: this.response,
              });
          } else
            triggerErrorEvent(
              getDocument().body,
              "htmx:historyCacheMissLoadError",
              Z,
            );
        }),
        Y.send();
    }
    function restoreHistory(J) {
      saveCurrentPageToHistory(),
        (J = J || location.pathname + location.search);
      let Y = getCachedHistory(J);
      if (Y) {
        let Z = makeFragment(Y.content),
          $ = getHistoryElement(),
          W = makeSettleInfo($);
        handleTitle(Y.title),
          handlePreservedElements(Z),
          swapInnerHTML($, Z, W),
          restorePreservedElements(),
          settleImmediately(W.tasks),
          getWindow().setTimeout(function () {
            window.scrollTo(0, Y.scroll);
          }, 0),
          (currentPathForHistory = J),
          triggerEvent(getDocument().body, "htmx:historyRestore", {
            path: J,
            item: Y,
          });
      } else if (htmx.config.refreshOnHistoryMiss) window.location.reload(!0);
      else loadHistoryFromServer(J);
    }
    function addRequestIndicatorClasses(J) {
      let Y = findAttributeTargets(J, "hx-indicator");
      if (Y == null) Y = [J];
      return (
        forEach(Y, function (Z) {
          let $ = getInternalData(Z);
          ($.requestCount = ($.requestCount || 0) + 1),
            Z.classList.add.call(Z.classList, htmx.config.requestClass);
        }),
        Y
      );
    }
    function disableElements(J) {
      let Y = findAttributeTargets(J, "hx-disabled-elt");
      if (Y == null) Y = [];
      return (
        forEach(Y, function (Z) {
          let $ = getInternalData(Z);
          ($.requestCount = ($.requestCount || 0) + 1),
            Z.setAttribute("disabled", ""),
            Z.setAttribute("data-disabled-by-htmx", "");
        }),
        Y
      );
    }
    function removeRequestIndicators(J, Y) {
      forEach(J.concat(Y), function (Z) {
        let $ = getInternalData(Z);
        $.requestCount = ($.requestCount || 1) - 1;
      }),
        forEach(J, function (Z) {
          if (getInternalData(Z).requestCount === 0)
            Z.classList.remove.call(Z.classList, htmx.config.requestClass);
        }),
        forEach(Y, function (Z) {
          if (getInternalData(Z).requestCount === 0)
            Z.removeAttribute("disabled"),
              Z.removeAttribute("data-disabled-by-htmx");
        });
    }
    function haveSeenNode(J, Y) {
      for (let Z = 0; Z < J.length; Z++) if (J[Z].isSameNode(Y)) return !0;
      return !1;
    }
    function shouldInclude(J) {
      let Y = J;
      if (
        Y.name === "" ||
        Y.name == null ||
        Y.disabled ||
        closest(Y, "fieldset[disabled]")
      )
        return !1;
      if (
        Y.type === "button" ||
        Y.type === "submit" ||
        Y.tagName === "image" ||
        Y.tagName === "reset" ||
        Y.tagName === "file"
      )
        return !1;
      if (Y.type === "checkbox" || Y.type === "radio") return Y.checked;
      return !0;
    }
    function addValueToFormData(J, Y, Z) {
      if (J != null && Y != null)
        if (Array.isArray(Y))
          Y.forEach(function ($) {
            Z.append(J, $);
          });
        else Z.append(J, Y);
    }
    function removeValueFromFormData(J, Y, Z) {
      if (J != null && Y != null) {
        let $ = Z.getAll(J);
        if (Array.isArray(Y)) $ = $.filter((W) => Y.indexOf(W) < 0);
        else $ = $.filter((W) => W !== Y);
        Z.delete(J), forEach($, (W) => Z.append(J, W));
      }
    }
    function processInputValue(J, Y, Z, $, W) {
      if ($ == null || haveSeenNode(J, $)) return;
      else J.push($);
      if (shouldInclude($)) {
        let X = getRawAttribute($, "name"),
          Q = $.value;
        if ($ instanceof HTMLSelectElement && $.multiple)
          Q = toArray($.querySelectorAll("option:checked")).map(function (G) {
            return G.value;
          });
        if ($ instanceof HTMLInputElement && $.files) Q = toArray($.files);
        if ((addValueToFormData(X, Q, Y), W)) validateElement($, Z);
      }
      if ($ instanceof HTMLFormElement)
        forEach($.elements, function (X) {
          if (J.indexOf(X) >= 0) removeValueFromFormData(X.name, X.value, Y);
          else J.push(X);
          if (W) validateElement(X, Z);
        }),
          new FormData($).forEach(function (X, Q) {
            if (X instanceof File && X.name === "") return;
            addValueToFormData(Q, X, Y);
          });
    }
    function validateElement(J, Y) {
      let Z = J;
      if (Z.willValidate) {
        if ((triggerEvent(Z, "htmx:validation:validate"), !Z.checkValidity()))
          Y.push({
            elt: Z,
            message: Z.validationMessage,
            validity: Z.validity,
          }),
            triggerEvent(Z, "htmx:validation:failed", {
              message: Z.validationMessage,
              validity: Z.validity,
            });
      }
    }
    function overrideFormData(J, Y) {
      for (let Z of Y.keys()) J.delete(Z);
      return (
        Y.forEach(function (Z, $) {
          J.append($, Z);
        }),
        J
      );
    }
    function getInputValues(J, Y) {
      let Z = [],
        $ = new FormData(),
        W = new FormData(),
        X = [],
        Q = getInternalData(J);
      if (Q.lastButtonClicked && !bodyContains(Q.lastButtonClicked))
        Q.lastButtonClicked = null;
      let G =
        (J instanceof HTMLFormElement && J.noValidate !== !0) ||
        getAttributeValue(J, "hx-validate") === "true";
      if (Q.lastButtonClicked)
        G = G && Q.lastButtonClicked.formNoValidate !== !0;
      if (Y !== "get") processInputValue(Z, W, X, closest(J, "form"), G);
      if (
        (processInputValue(Z, $, X, J, G),
        Q.lastButtonClicked ||
          J.tagName === "BUTTON" ||
          (J.tagName === "INPUT" && getRawAttribute(J, "type") === "submit"))
      ) {
        let U = Q.lastButtonClicked || J,
          _ = getRawAttribute(U, "name");
        addValueToFormData(_, U.value, W);
      }
      let B = findAttributeTargets(J, "hx-include");
      return (
        forEach(B, function (U) {
          if (
            (processInputValue(Z, $, X, asElement(U), G), !matches(U, "form"))
          )
            forEach(
              asParentNode(U).querySelectorAll(INPUT_SELECTOR),
              function (_) {
                processInputValue(Z, $, X, _, G);
              },
            );
        }),
        overrideFormData($, W),
        { errors: X, formData: $, values: formDataProxy($) }
      );
    }
    function appendParam(J, Y, Z) {
      if (J !== "") J += "&";
      if (String(Z) === "[object Object]") Z = JSON.stringify(Z);
      let $ = encodeURIComponent(Z);
      return (J += encodeURIComponent(Y) + "=" + $), J;
    }
    function urlEncode(J) {
      J = formDataFromObject(J);
      let Y = "";
      return (
        J.forEach(function (Z, $) {
          Y = appendParam(Y, $, Z);
        }),
        Y
      );
    }
    function getHeaders(J, Y, Z) {
      let $ = {
        "HX-Request": "true",
        "HX-Trigger": getRawAttribute(J, "id"),
        "HX-Trigger-Name": getRawAttribute(J, "name"),
        "HX-Target": getAttributeValue(Y, "id"),
        "HX-Current-URL": getDocument().location.href,
      };
      if ((getValuesForElement(J, "hx-headers", !1, $), Z !== void 0))
        $["HX-Prompt"] = Z;
      if (getInternalData(J).boosted) $["HX-Boosted"] = "true";
      return $;
    }
    function filterValues(J, Y) {
      let Z = getClosestAttributeValue(Y, "hx-params");
      if (Z)
        if (Z === "none") return new FormData();
        else if (Z === "*") return J;
        else if (Z.indexOf("not ") === 0)
          return (
            forEach(Z.slice(4).split(","), function ($) {
              ($ = $.trim()), J.delete($);
            }),
            J
          );
        else {
          let $ = new FormData();
          return (
            forEach(Z.split(","), function (W) {
              if (((W = W.trim()), J.has(W)))
                J.getAll(W).forEach(function (X) {
                  $.append(W, X);
                });
            }),
            $
          );
        }
      else return J;
    }
    function isAnchorLink(J) {
      return (
        !!getRawAttribute(J, "href") &&
        getRawAttribute(J, "href").indexOf("#") >= 0
      );
    }
    function getSwapSpecification(J, Y) {
      let Z = Y || getClosestAttributeValue(J, "hx-swap"),
        $ = {
          swapStyle: getInternalData(J).boosted
            ? "innerHTML"
            : htmx.config.defaultSwapStyle,
          swapDelay: htmx.config.defaultSwapDelay,
          settleDelay: htmx.config.defaultSettleDelay,
        };
      if (
        htmx.config.scrollIntoViewOnBoost &&
        getInternalData(J).boosted &&
        !isAnchorLink(J)
      )
        $.show = "top";
      if (Z) {
        let Q = splitOnWhitespace(Z);
        if (Q.length > 0)
          for (let G = 0; G < Q.length; G++) {
            let B = Q[G];
            if (B.indexOf("swap:") === 0)
              $.swapDelay = parseInterval(B.slice(5));
            else if (B.indexOf("settle:") === 0)
              $.settleDelay = parseInterval(B.slice(7));
            else if (B.indexOf("transition:") === 0)
              $.transition = B.slice(11) === "true";
            else if (B.indexOf("ignoreTitle:") === 0)
              $.ignoreTitle = B.slice(12) === "true";
            else if (B.indexOf("scroll:") === 0) {
              var W = B.slice(7).split(":");
              let _ = W.pop();
              var X = W.length > 0 ? W.join(":") : null;
              ($.scroll = _), ($.scrollTarget = X);
            } else if (B.indexOf("show:") === 0) {
              var W = B.slice(5).split(":");
              let z = W.pop();
              var X = W.length > 0 ? W.join(":") : null;
              ($.show = z), ($.showTarget = X);
            } else if (B.indexOf("focus-scroll:") === 0) {
              let U = B.slice(13);
              $.focusScroll = U == "true";
            } else if (G == 0) $.swapStyle = B;
            else logError("Unknown modifier in hx-swap: " + B);
          }
      }
      return $;
    }
    function usesFormData(J) {
      return (
        getClosestAttributeValue(J, "hx-encoding") === "multipart/form-data" ||
        (matches(J, "form") &&
          getRawAttribute(J, "enctype") === "multipart/form-data")
      );
    }
    function encodeParamsForBody(J, Y, Z) {
      let $ = null;
      if (
        (withExtensions(Y, function (W) {
          if ($ == null) $ = W.encodeParameters(J, Z, Y);
        }),
        $ != null)
      )
        return $;
      else if (usesFormData(Y))
        return overrideFormData(new FormData(), formDataFromObject(Z));
      else return urlEncode(Z);
    }
    function makeSettleInfo(J) {
      return { tasks: [], elts: [J] };
    }
    function updateScrollState(J, Y) {
      let Z = J[0],
        $ = J[J.length - 1];
      if (Y.scroll) {
        var W = null;
        if (Y.scrollTarget) W = asElement(querySelectorExt(Z, Y.scrollTarget));
        if (Y.scroll === "top" && (Z || W)) (W = W || Z), (W.scrollTop = 0);
        if (Y.scroll === "bottom" && ($ || W))
          (W = W || $), (W.scrollTop = W.scrollHeight);
      }
      if (Y.show) {
        var W = null;
        if (Y.showTarget) {
          let Q = Y.showTarget;
          if (Y.showTarget === "window") Q = "body";
          W = asElement(querySelectorExt(Z, Q));
        }
        if (Y.show === "top" && (Z || W))
          (W = W || Z),
            W.scrollIntoView({
              block: "start",
              behavior: htmx.config.scrollBehavior,
            });
        if (Y.show === "bottom" && ($ || W))
          (W = W || $),
            W.scrollIntoView({
              block: "end",
              behavior: htmx.config.scrollBehavior,
            });
      }
    }
    function getValuesForElement(J, Y, Z, $) {
      if ($ == null) $ = {};
      if (J == null) return $;
      let W = getAttributeValue(J, Y);
      if (W) {
        let X = W.trim(),
          Q = Z;
        if (X === "unset") return null;
        if (X.indexOf("javascript:") === 0) (X = X.slice(11)), (Q = !0);
        else if (X.indexOf("js:") === 0) (X = X.slice(3)), (Q = !0);
        if (X.indexOf("{") !== 0) X = "{" + X + "}";
        let G;
        if (Q)
          G = maybeEval(
            J,
            function () {
              return Function("return (" + X + ")")();
            },
            {},
          );
        else G = parseJSON(X);
        for (let B in G)
          if (G.hasOwnProperty(B)) {
            if ($[B] == null) $[B] = G[B];
          }
      }
      return getValuesForElement(asElement(parentElt(J)), Y, Z, $);
    }
    function maybeEval(J, Y, Z) {
      if (htmx.config.allowEval) return Y();
      else return triggerErrorEvent(J, "htmx:evalDisallowedError"), Z;
    }
    function getHXVarsForElement(J, Y) {
      return getValuesForElement(J, "hx-vars", !0, Y);
    }
    function getHXValsForElement(J, Y) {
      return getValuesForElement(J, "hx-vals", !1, Y);
    }
    function getExpressionVars(J) {
      return mergeObjects(getHXVarsForElement(J), getHXValsForElement(J));
    }
    function safelySetHeaderValue(J, Y, Z) {
      if (Z !== null)
        try {
          J.setRequestHeader(Y, Z);
        } catch ($) {
          J.setRequestHeader(Y, encodeURIComponent(Z)),
            J.setRequestHeader(Y + "-URI-AutoEncoded", "true");
        }
    }
    function getPathFromResponse(J) {
      if (J.responseURL && typeof URL !== "undefined")
        try {
          let Y = new URL(J.responseURL);
          return Y.pathname + Y.search;
        } catch (Y) {
          triggerErrorEvent(getDocument().body, "htmx:badResponseUrl", {
            url: J.responseURL,
          });
        }
    }
    function hasHeader(J, Y) {
      return Y.test(J.getAllResponseHeaders());
    }
    function ajaxHelper(J, Y, Z) {
      if (((J = J.toLowerCase()), Z))
        if (Z instanceof Element || typeof Z === "string")
          return issueAjaxRequest(J, Y, null, null, {
            targetOverride: resolveTarget(Z) || DUMMY_ELT,
            returnPromise: !0,
          });
        else {
          let $ = resolveTarget(Z.target);
          if ((Z.target && !$) || (Z.source && !$ && !resolveTarget(Z.source)))
            $ = DUMMY_ELT;
          return issueAjaxRequest(J, Y, resolveTarget(Z.source), Z.event, {
            handler: Z.handler,
            headers: Z.headers,
            values: Z.values,
            targetOverride: $,
            swapOverride: Z.swap,
            select: Z.select,
            returnPromise: !0,
          });
        }
      else return issueAjaxRequest(J, Y, null, null, { returnPromise: !0 });
    }
    function hierarchyForElt(J) {
      let Y = [];
      while (J) Y.push(J), (J = J.parentElement);
      return Y;
    }
    function verifyPath(J, Y, Z) {
      let $, W;
      if (typeof URL === "function")
        (W = new URL(Y, document.location.href)),
          ($ = document.location.origin === W.origin);
      else (W = Y), ($ = startsWith(Y, document.location.origin));
      if (htmx.config.selfRequestsOnly) {
        if (!$) return !1;
      }
      return triggerEvent(
        J,
        "htmx:validateUrl",
        mergeObjects({ url: W, sameHost: $ }, Z),
      );
    }
    function formDataFromObject(J) {
      if (J instanceof FormData) return J;
      let Y = new FormData();
      for (let Z in J)
        if (J.hasOwnProperty(Z))
          if (J[Z] && typeof J[Z].forEach === "function")
            J[Z].forEach(function ($) {
              Y.append(Z, $);
            });
          else if (typeof J[Z] === "object" && !(J[Z] instanceof Blob))
            Y.append(Z, JSON.stringify(J[Z]));
          else Y.append(Z, J[Z]);
      return Y;
    }
    function formDataArrayProxy(J, Y, Z) {
      return new Proxy(Z, {
        get: function ($, W) {
          if (typeof W === "number") return $[W];
          if (W === "length") return $.length;
          if (W === "push")
            return function (X) {
              $.push(X), J.append(Y, X);
            };
          if (typeof $[W] === "function")
            return function () {
              $[W].apply($, arguments),
                J.delete(Y),
                $.forEach(function (X) {
                  J.append(Y, X);
                });
            };
          if ($[W] && $[W].length === 1) return $[W][0];
          else return $[W];
        },
        set: function ($, W, X) {
          return (
            ($[W] = X),
            J.delete(Y),
            $.forEach(function (Q) {
              J.append(Y, Q);
            }),
            !0
          );
        },
      });
    }
    function formDataProxy(J) {
      return new Proxy(J, {
        get: function (Y, Z) {
          if (typeof Z === "symbol") {
            let W = Reflect.get(Y, Z);
            if (typeof W === "function")
              return function () {
                return W.apply(J, arguments);
              };
            else return W;
          }
          if (Z === "toJSON") return () => Object.fromEntries(J);
          if (Z in Y)
            if (typeof Y[Z] === "function")
              return function () {
                return J[Z].apply(J, arguments);
              };
            else return Y[Z];
          let $ = J.getAll(Z);
          if ($.length === 0) return;
          else if ($.length === 1) return $[0];
          else return formDataArrayProxy(Y, Z, $);
        },
        set: function (Y, Z, $) {
          if (typeof Z !== "string") return !1;
          if ((Y.delete(Z), $ && typeof $.forEach === "function"))
            $.forEach(function (W) {
              Y.append(Z, W);
            });
          else if (typeof $ === "object" && !($ instanceof Blob))
            Y.append(Z, JSON.stringify($));
          else Y.append(Z, $);
          return !0;
        },
        deleteProperty: function (Y, Z) {
          if (typeof Z === "string") Y.delete(Z);
          return !0;
        },
        ownKeys: function (Y) {
          return Reflect.ownKeys(Object.fromEntries(Y));
        },
        getOwnPropertyDescriptor: function (Y, Z) {
          return Reflect.getOwnPropertyDescriptor(Object.fromEntries(Y), Z);
        },
      });
    }
    function issueAjaxRequest(J, Y, Z, $, W, X) {
      let Q = null,
        G = null;
      if (
        ((W = W != null ? W : {}),
        W.returnPromise && typeof Promise !== "undefined")
      )
        var B = new Promise(function (F, V) {
          (Q = F), (G = V);
        });
      if (Z == null) Z = getDocument().body;
      let U = W.handler || handleAjaxResponse,
        _ = W.select || null;
      if (!bodyContains(Z)) return maybeCall(Q), B;
      let z = W.targetOverride || asElement(getTarget(Z));
      if (z == null || z == DUMMY_ELT)
        return (
          triggerErrorEvent(Z, "htmx:targetError", {
            target: getAttributeValue(Z, "hx-target"),
          }),
          maybeCall(G),
          B
        );
      let K = getInternalData(Z),
        q = K.lastButtonClicked;
      if (q) {
        let F = getRawAttribute(q, "formaction");
        if (F != null) Y = F;
        let V = getRawAttribute(q, "formmethod");
        if (V != null) {
          if (V.toLowerCase() !== "dialog") J = V;
        }
      }
      let j = getClosestAttributeValue(Z, "hx-confirm");
      if (X === void 0) {
        if (
          triggerEvent(Z, "htmx:confirm", {
            target: z,
            elt: Z,
            path: Y,
            verb: J,
            triggeringEvent: $,
            etc: W,
            issueRequest: function (v) {
              return issueAjaxRequest(J, Y, Z, $, W, !!v);
            },
            question: j,
          }) === !1
        )
          return maybeCall(Q), B;
      }
      let T = Z,
        M = getClosestAttributeValue(Z, "hx-sync"),
        L = null,
        P = !1;
      if (M) {
        let F = M.split(":"),
          V = F[0].trim();
        if (V === "this") T = findThisElement(Z, "hx-sync");
        else T = asElement(querySelectorExt(Z, V));
        if (
          ((M = (F[1] || "drop").trim()),
          (K = getInternalData(T)),
          M === "drop" && K.xhr && K.abortable !== !0)
        )
          return maybeCall(Q), B;
        else if (M === "abort")
          if (K.xhr) return maybeCall(Q), B;
          else P = !0;
        else if (M === "replace") triggerEvent(T, "htmx:abort");
        else if (M.indexOf("queue") === 0)
          L = (M.split(" ")[1] || "last").trim();
      }
      if (K.xhr)
        if (K.abortable) triggerEvent(T, "htmx:abort");
        else {
          if (L == null) {
            if ($) {
              let F = getInternalData($);
              if (F && F.triggerSpec && F.triggerSpec.queue)
                L = F.triggerSpec.queue;
            }
            if (L == null) L = "last";
          }
          if (K.queuedRequests == null) K.queuedRequests = [];
          if (L === "first" && K.queuedRequests.length === 0)
            K.queuedRequests.push(function () {
              issueAjaxRequest(J, Y, Z, $, W);
            });
          else if (L === "all")
            K.queuedRequests.push(function () {
              issueAjaxRequest(J, Y, Z, $, W);
            });
          else if (L === "last")
            (K.queuedRequests = []),
              K.queuedRequests.push(function () {
                issueAjaxRequest(J, Y, Z, $, W);
              });
          return maybeCall(Q), B;
        }
      let A = new XMLHttpRequest();
      (K.xhr = A), (K.abortable = P);
      let H = function () {
          if (
            ((K.xhr = null),
            (K.abortable = !1),
            K.queuedRequests != null && K.queuedRequests.length > 0)
          )
            K.queuedRequests.shift()();
        },
        x = getClosestAttributeValue(Z, "hx-prompt");
      if (x) {
        var D = prompt(x);
        if (
          D === null ||
          !triggerEvent(Z, "htmx:prompt", { prompt: D, target: z })
        )
          return maybeCall(Q), H(), B;
      }
      if (j && !X) {
        if (!confirm(j)) return maybeCall(Q), H(), B;
      }
      let w = getHeaders(Z, z, D);
      if (J !== "get" && !usesFormData(Z))
        w["Content-Type"] = "application/x-www-form-urlencoded";
      if (W.headers) w = mergeObjects(w, W.headers);
      let E = getInputValues(Z, J),
        d = E.errors,
        JJ = E.formData;
      if (W.values) overrideFormData(JJ, formDataFromObject(W.values));
      let gJ = formDataFromObject(getExpressionVars(Z)),
        mJ = overrideFormData(JJ, gJ),
        i = filterValues(mJ, Z);
      if (htmx.config.getCacheBusterParam && J === "get")
        i.set("org.htmx.cache-buster", getRawAttribute(z, "id") || "true");
      if (Y == null || Y === "") Y = getDocument().location.href;
      let cJ = getValuesForElement(Z, "hx-request"),
        SY = getInternalData(Z).boosted,
        RJ = htmx.config.methodsThatUseUrlParams.indexOf(J) >= 0,
        b = {
          boosted: SY,
          useUrlParams: RJ,
          formData: i,
          parameters: formDataProxy(i),
          unfilteredFormData: mJ,
          unfilteredParameters: formDataProxy(mJ),
          headers: w,
          target: z,
          verb: J,
          errors: d,
          withCredentials:
            W.credentials || cJ.credentials || htmx.config.withCredentials,
          timeout: W.timeout || cJ.timeout || htmx.config.timeout,
          path: Y,
          triggeringEvent: $,
        };
      if (!triggerEvent(Z, "htmx:configRequest", b))
        return maybeCall(Q), H(), B;
      if (
        ((Y = b.path),
        (J = b.verb),
        (w = b.headers),
        (i = formDataFromObject(b.parameters)),
        (d = b.errors),
        (RJ = b.useUrlParams),
        d && d.length > 0)
      )
        return (
          triggerEvent(Z, "htmx:validation:halted", b), maybeCall(Q), H(), B
        );
      let yY = Y.split("#"),
        A0 = yY[0],
        pJ = yY[1],
        f = Y;
      if (RJ) {
        if (((f = A0), !i.keys().next().done)) {
          if (f.indexOf("?") < 0) f += "?";
          else f += "&";
          if (((f += urlEncode(i)), pJ)) f += "#" + pJ;
        }
      }
      if (!verifyPath(Z, f, b))
        return triggerErrorEvent(Z, "htmx:invalidPath", b), maybeCall(G), B;
      if (
        (A.open(J.toUpperCase(), f, !0),
        A.overrideMimeType("text/html"),
        (A.withCredentials = b.withCredentials),
        (A.timeout = b.timeout),
        cJ.noHeaders)
      );
      else
        for (let F in w)
          if (w.hasOwnProperty(F)) {
            let V = w[F];
            safelySetHeaderValue(A, F, V);
          }
      let N = {
        xhr: A,
        target: z,
        requestConfig: b,
        etc: W,
        boosted: SY,
        select: _,
        pathInfo: {
          requestPath: Y,
          finalRequestPath: f,
          responsePath: null,
          anchor: pJ,
        },
      };
      if (
        ((A.onload = function () {
          try {
            let F = hierarchyForElt(Z);
            if (
              ((N.pathInfo.responsePath = getPathFromResponse(A)),
              U(Z, N),
              N.keepIndicators !== !0)
            )
              removeRequestIndicators(HJ, CJ);
            if (
              (triggerEvent(Z, "htmx:afterRequest", N),
              triggerEvent(Z, "htmx:afterOnLoad", N),
              !bodyContains(Z))
            ) {
              let V = null;
              while (F.length > 0 && V == null) {
                let v = F.shift();
                if (bodyContains(v)) V = v;
              }
              if (V)
                triggerEvent(V, "htmx:afterRequest", N),
                  triggerEvent(V, "htmx:afterOnLoad", N);
            }
            maybeCall(Q), H();
          } catch (F) {
            throw (
              (triggerErrorEvent(
                Z,
                "htmx:onLoadError",
                mergeObjects({ error: F }, N),
              ),
              F)
            );
          }
        }),
        (A.onerror = function () {
          removeRequestIndicators(HJ, CJ),
            triggerErrorEvent(Z, "htmx:afterRequest", N),
            triggerErrorEvent(Z, "htmx:sendError", N),
            maybeCall(G),
            H();
        }),
        (A.onabort = function () {
          removeRequestIndicators(HJ, CJ),
            triggerErrorEvent(Z, "htmx:afterRequest", N),
            triggerErrorEvent(Z, "htmx:sendAbort", N),
            maybeCall(G),
            H();
        }),
        (A.ontimeout = function () {
          removeRequestIndicators(HJ, CJ),
            triggerErrorEvent(Z, "htmx:afterRequest", N),
            triggerErrorEvent(Z, "htmx:timeout", N),
            maybeCall(G),
            H();
        }),
        !triggerEvent(Z, "htmx:beforeRequest", N))
      )
        return maybeCall(Q), H(), B;
      var HJ = addRequestIndicatorClasses(Z),
        CJ = disableElements(Z);
      forEach(["loadstart", "loadend", "progress", "abort"], function (F) {
        forEach([A, A.upload], function (V) {
          V.addEventListener(F, function (v) {
            triggerEvent(Z, "htmx:xhr:" + F, {
              lengthComputable: v.lengthComputable,
              loaded: v.loaded,
              total: v.total,
            });
          });
        });
      }),
        triggerEvent(Z, "htmx:beforeSend", N);
      let P0 = RJ ? null : encodeParamsForBody(A, Z, i);
      return A.send(P0), B;
    }
    function determineHistoryUpdates(J, Y) {
      let Z = Y.xhr,
        $ = null,
        W = null;
      if (hasHeader(Z, /HX-Push:/i))
        ($ = Z.getResponseHeader("HX-Push")), (W = "push");
      else if (hasHeader(Z, /HX-Push-Url:/i))
        ($ = Z.getResponseHeader("HX-Push-Url")), (W = "push");
      else if (hasHeader(Z, /HX-Replace-Url:/i))
        ($ = Z.getResponseHeader("HX-Replace-Url")), (W = "replace");
      if ($)
        if ($ === "false") return {};
        else return { type: W, path: $ };
      let X = Y.pathInfo.finalRequestPath,
        Q = Y.pathInfo.responsePath,
        G = getClosestAttributeValue(J, "hx-push-url"),
        B = getClosestAttributeValue(J, "hx-replace-url"),
        U = getInternalData(J).boosted,
        _ = null,
        z = null;
      if (G) (_ = "push"), (z = G);
      else if (B) (_ = "replace"), (z = B);
      else if (U) (_ = "push"), (z = Q || X);
      if (z) {
        if (z === "false") return {};
        if (z === "true") z = Q || X;
        if (Y.pathInfo.anchor && z.indexOf("#") === -1)
          z = z + "#" + Y.pathInfo.anchor;
        return { type: _, path: z };
      } else return {};
    }
    function codeMatches(J, Y) {
      var Z = new RegExp(J.code);
      return Z.test(Y.toString(10));
    }
    function resolveResponseHandling(J) {
      for (var Y = 0; Y < htmx.config.responseHandling.length; Y++) {
        var Z = htmx.config.responseHandling[Y];
        if (codeMatches(Z, J.status)) return Z;
      }
      return { swap: !1 };
    }
    function handleTitle(J) {
      if (J) {
        let Y = find("title");
        if (Y) Y.innerHTML = J;
        else window.document.title = J;
      }
    }
    function handleAjaxResponse(J, Y) {
      let { xhr: Z, target: $, etc: W, select: X } = Y;
      if (!triggerEvent(J, "htmx:beforeOnLoad", Y)) return;
      if (hasHeader(Z, /HX-Trigger:/i)) handleTriggerHeader(Z, "HX-Trigger", J);
      if (hasHeader(Z, /HX-Location:/i)) {
        saveCurrentPageToHistory();
        let P = Z.getResponseHeader("HX-Location");
        var Q;
        if (P.indexOf("{") === 0)
          (Q = parseJSON(P)), (P = Q.path), delete Q.path;
        ajaxHelper("get", P, Q).then(function () {
          pushUrlIntoHistory(P);
        });
        return;
      }
      let G =
        hasHeader(Z, /HX-Refresh:/i) &&
        Z.getResponseHeader("HX-Refresh") === "true";
      if (hasHeader(Z, /HX-Redirect:/i)) {
        (Y.keepIndicators = !0),
          (location.href = Z.getResponseHeader("HX-Redirect")),
          G && location.reload();
        return;
      }
      if (G) {
        (Y.keepIndicators = !0), location.reload();
        return;
      }
      if (hasHeader(Z, /HX-Retarget:/i))
        if (Z.getResponseHeader("HX-Retarget") === "this") Y.target = J;
        else
          Y.target = asElement(
            querySelectorExt(J, Z.getResponseHeader("HX-Retarget")),
          );
      let B = determineHistoryUpdates(J, Y),
        U = resolveResponseHandling(Z),
        _ = U.swap,
        z = !!U.error,
        K = htmx.config.ignoreTitle || U.ignoreTitle,
        q = U.select;
      if (U.target) Y.target = asElement(querySelectorExt(J, U.target));
      var j = W.swapOverride;
      if (j == null && U.swapOverride) j = U.swapOverride;
      if (hasHeader(Z, /HX-Retarget:/i))
        if (Z.getResponseHeader("HX-Retarget") === "this") Y.target = J;
        else
          Y.target = asElement(
            querySelectorExt(J, Z.getResponseHeader("HX-Retarget")),
          );
      if (hasHeader(Z, /HX-Reswap:/i)) j = Z.getResponseHeader("HX-Reswap");
      var T = Z.response,
        M = mergeObjects(
          {
            shouldSwap: _,
            serverResponse: T,
            isError: z,
            ignoreTitle: K,
            selectOverride: q,
            swapOverride: j,
          },
          Y,
        );
      if (U.event && !triggerEvent($, U.event, M)) return;
      if (!triggerEvent($, "htmx:beforeSwap", M)) return;
      if (
        (($ = M.target),
        (T = M.serverResponse),
        (z = M.isError),
        (K = M.ignoreTitle),
        (q = M.selectOverride),
        (j = M.swapOverride),
        (Y.target = $),
        (Y.failed = z),
        (Y.successful = !z),
        M.shouldSwap)
      ) {
        if (Z.status === 286) cancelPolling(J);
        if (
          (withExtensions(J, function (E) {
            T = E.transformResponse(T, Z, J);
          }),
          B.type)
        )
          saveCurrentPageToHistory();
        var L = getSwapSpecification(J, j);
        if (!L.hasOwnProperty("ignoreTitle")) L.ignoreTitle = K;
        $.classList.add(htmx.config.swappingClass);
        let P = null,
          A = null;
        if (X) q = X;
        if (hasHeader(Z, /HX-Reselect:/i))
          q = Z.getResponseHeader("HX-Reselect");
        let H = getClosestAttributeValue(J, "hx-select-oob"),
          x = getClosestAttributeValue(J, "hx-select"),
          D = function () {
            try {
              if (B.type)
                if (
                  (triggerEvent(
                    getDocument().body,
                    "htmx:beforeHistoryUpdate",
                    mergeObjects({ history: B }, Y),
                  ),
                  B.type === "push")
                )
                  pushUrlIntoHistory(B.path),
                    triggerEvent(getDocument().body, "htmx:pushedIntoHistory", {
                      path: B.path,
                    });
                else
                  replaceUrlInHistory(B.path),
                    triggerEvent(getDocument().body, "htmx:replacedInHistory", {
                      path: B.path,
                    });
              swap($, T, L, {
                select: q || x,
                selectOOB: H,
                eventInfo: Y,
                anchor: Y.pathInfo.anchor,
                contextElement: J,
                afterSwapCallback: function () {
                  if (hasHeader(Z, /HX-Trigger-After-Swap:/i)) {
                    let E = J;
                    if (!bodyContains(J)) E = getDocument().body;
                    handleTriggerHeader(Z, "HX-Trigger-After-Swap", E);
                  }
                },
                afterSettleCallback: function () {
                  if (hasHeader(Z, /HX-Trigger-After-Settle:/i)) {
                    let E = J;
                    if (!bodyContains(J)) E = getDocument().body;
                    handleTriggerHeader(Z, "HX-Trigger-After-Settle", E);
                  }
                  maybeCall(P);
                },
              });
            } catch (E) {
              throw (
                (triggerErrorEvent(J, "htmx:swapError", Y), maybeCall(A), E)
              );
            }
          },
          w = htmx.config.globalViewTransitions;
        if (L.hasOwnProperty("transition")) w = L.transition;
        if (
          w &&
          triggerEvent(J, "htmx:beforeTransition", Y) &&
          typeof Promise !== "undefined" &&
          document.startViewTransition
        ) {
          let E = new Promise(function (JJ, gJ) {
              (P = JJ), (A = gJ);
            }),
            d = D;
          D = function () {
            document.startViewTransition(function () {
              return d(), E;
            });
          };
        }
        if (L.swapDelay > 0) getWindow().setTimeout(D, L.swapDelay);
        else D();
      }
      if (z)
        triggerErrorEvent(
          J,
          "htmx:responseError",
          mergeObjects(
            {
              error:
                "Response Status Error Code " +
                Z.status +
                " from " +
                Y.pathInfo.requestPath,
            },
            Y,
          ),
        );
    }
    let extensions = {};
    function extensionBase() {
      return {
        init: function (J) {
          return null;
        },
        getSelectors: function () {
          return null;
        },
        onEvent: function (J, Y) {
          return !0;
        },
        transformResponse: function (J, Y, Z) {
          return J;
        },
        isInlineSwap: function (J) {
          return !1;
        },
        handleSwap: function (J, Y, Z, $) {
          return !1;
        },
        encodeParameters: function (J, Y, Z) {
          return null;
        },
      };
    }
    function defineExtension(J, Y) {
      if (Y.init) Y.init(internalAPI);
      extensions[J] = mergeObjects(extensionBase(), Y);
    }
    function removeExtension(J) {
      delete extensions[J];
    }
    function getExtensions(J, Y, Z) {
      if (Y == null) Y = [];
      if (J == null) return Y;
      if (Z == null) Z = [];
      let $ = getAttributeValue(J, "hx-ext");
      if ($)
        forEach($.split(","), function (W) {
          if (((W = W.replace(/ /g, "")), W.slice(0, 7) == "ignore:")) {
            Z.push(W.slice(7));
            return;
          }
          if (Z.indexOf(W) < 0) {
            let X = extensions[W];
            if (X && Y.indexOf(X) < 0) Y.push(X);
          }
        });
      return getExtensions(asElement(parentElt(J)), Y, Z);
    }
    var isReady = !1;
    getDocument().addEventListener("DOMContentLoaded", function () {
      isReady = !0;
    });
    function ready(J) {
      if (isReady || getDocument().readyState === "complete") J();
      else getDocument().addEventListener("DOMContentLoaded", J);
    }
    function insertIndicatorStyles() {
      if (htmx.config.includeIndicatorStyles !== !1) {
        let J = htmx.config.inlineStyleNonce
          ? ` nonce="${htmx.config.inlineStyleNonce}"`
          : "";
        getDocument().head.insertAdjacentHTML(
          "beforeend",
          "<style" +
            J +
            ">      ." +
            htmx.config.indicatorClass +
            "{opacity:0}      ." +
            htmx.config.requestClass +
            " ." +
            htmx.config.indicatorClass +
            "{opacity:1; transition: opacity 200ms ease-in;}      ." +
            htmx.config.requestClass +
            "." +
            htmx.config.indicatorClass +
            "{opacity:1; transition: opacity 200ms ease-in;}      </style>",
        );
      }
    }
    function getMetaConfig() {
      let J = getDocument().querySelector('meta[name="htmx-config"]');
      if (J) return parseJSON(J.content);
      else return null;
    }
    function mergeMetaConfig() {
      let J = getMetaConfig();
      if (J) htmx.config = mergeObjects(htmx.config, J);
    }
    return (
      ready(function () {
        mergeMetaConfig(), insertIndicatorStyles();
        let J = getDocument().body;
        processNode(J);
        let Y = getDocument().querySelectorAll(
          "[hx-trigger='restored'],[data-hx-trigger='restored']",
        );
        J.addEventListener("htmx:abort", function ($) {
          let W = $.target,
            X = getInternalData(W);
          if (X && X.xhr) X.xhr.abort();
        });
        let Z = window.onpopstate ? window.onpopstate.bind(window) : null;
        (window.onpopstate = function ($) {
          if ($.state && $.state.htmx)
            restoreHistory(),
              forEach(Y, function (W) {
                triggerEvent(W, "htmx:restored", {
                  document: getDocument(),
                  triggerEvent,
                });
              });
          else if (Z) Z($);
        }),
          getWindow().setTimeout(function () {
            triggerEvent(J, "htmx:load", {}), (J = null);
          }, 0);
      }),
      htmx
    );
  })(),
  hY = F0;
var aJ = !1,
  lJ = !1,
  s = [],
  tJ = -1;
function R0(J) {
  H0(J);
}
function H0(J) {
  if (!s.includes(J)) s.push(J);
  j0();
}
function C0(J) {
  let Y = s.indexOf(J);
  if (Y !== -1 && Y > tJ) s.splice(Y, 1);
}
function j0() {
  if (!lJ && !aJ) (aJ = !0), queueMicrotask(O0);
}
function O0() {
  (aJ = !1), (lJ = !0);
  for (let J = 0; J < s.length; J++) s[J](), (tJ = J);
  (s.length = 0), (tJ = -1), (lJ = !1);
}
var $J,
  e,
  WJ,
  lY,
  eJ = !0;
function T0(J) {
  (eJ = !1), J(), (eJ = !0);
}
function V0(J) {
  ($J = J.reactive),
    (WJ = J.release),
    (e = (Y) =>
      J.effect(Y, {
        scheduler: (Z) => {
          if (eJ) R0(Z);
          else Z();
        },
      })),
    (lY = J.raw);
}
function fY(J) {
  e = J;
}
function E0(J) {
  let Y = () => {};
  return [
    ($) => {
      let W = e($);
      if (!J._x_effects)
        (J._x_effects = new Set()),
          (J._x_runEffects = () => {
            J._x_effects.forEach((X) => X());
          });
      return (
        J._x_effects.add(W),
        (Y = () => {
          if (W === void 0) return;
          J._x_effects.delete(W), WJ(W);
        }),
        W
      );
    },
    () => {
      Y();
    },
  ];
}
function tY(J, Y) {
  let Z = !0,
    $,
    W = e(() => {
      let X = J();
      if ((JSON.stringify(X), !Z))
        queueMicrotask(() => {
          Y(X, $), ($ = X);
        });
      else $ = X;
      Z = !1;
    });
  return () => WJ(W);
}
var eY = [],
  JZ = [],
  YZ = [];
function N0(J) {
  YZ.push(J);
}
function MY(J, Y) {
  if (typeof Y === "function") {
    if (!J._x_cleanups) J._x_cleanups = [];
    J._x_cleanups.push(Y);
  } else (Y = J), JZ.push(Y);
}
function ZZ(J) {
  eY.push(J);
}
function $Z(J, Y, Z) {
  if (!J._x_attributeCleanups) J._x_attributeCleanups = {};
  if (!J._x_attributeCleanups[Y]) J._x_attributeCleanups[Y] = [];
  J._x_attributeCleanups[Y].push(Z);
}
function WZ(J, Y) {
  if (!J._x_attributeCleanups) return;
  Object.entries(J._x_attributeCleanups).forEach(([Z, $]) => {
    if (Y === void 0 || Y.includes(Z))
      $.forEach((W) => W()), delete J._x_attributeCleanups[Z];
  });
}
function I0(J) {
  J._x_effects?.forEach(C0);
  while (J._x_cleanups?.length) J._x_cleanups.pop()();
}
var LY = new MutationObserver(FY),
  qY = !1;
function AY() {
  LY.observe(document, {
    subtree: !0,
    childList: !0,
    attributes: !0,
    attributeOldValue: !0,
  }),
    (qY = !0);
}
function XZ() {
  D0(), LY.disconnect(), (qY = !1);
}
var BJ = [];
function D0() {
  let J = LY.takeRecords();
  BJ.push(() => J.length > 0 && FY(J));
  let Y = BJ.length;
  queueMicrotask(() => {
    if (BJ.length === Y) while (BJ.length > 0) BJ.shift()();
  });
}
function C(J) {
  if (!qY) return J();
  XZ();
  let Y = J();
  return AY(), Y;
}
var PY = !1,
  DJ = [];
function w0() {
  PY = !0;
}
function b0() {
  (PY = !1), FY(DJ), (DJ = []);
}
function FY(J) {
  if (PY) {
    DJ = DJ.concat(J);
    return;
  }
  let Y = [],
    Z = new Set(),
    $ = new Map(),
    W = new Map();
  for (let X = 0; X < J.length; X++) {
    if (J[X].target._x_ignoreMutationObserver) continue;
    if (J[X].type === "childList")
      J[X].removedNodes.forEach((Q) => {
        if (Q.nodeType !== 1) return;
        if (!Q._x_marker) return;
        Z.add(Q);
      }),
        J[X].addedNodes.forEach((Q) => {
          if (Q.nodeType !== 1) return;
          if (Z.has(Q)) {
            Z.delete(Q);
            return;
          }
          if (Q._x_marker) return;
          Y.push(Q);
        });
    if (J[X].type === "attributes") {
      let Q = J[X].target,
        G = J[X].attributeName,
        B = J[X].oldValue,
        U = () => {
          if (!$.has(Q)) $.set(Q, []);
          $.get(Q).push({ name: G, value: Q.getAttribute(G) });
        },
        _ = () => {
          if (!W.has(Q)) W.set(Q, []);
          W.get(Q).push(G);
        };
      if (Q.hasAttribute(G) && B === null) U();
      else if (Q.hasAttribute(G)) _(), U();
      else _();
    }
  }
  W.forEach((X, Q) => {
    WZ(Q, X);
  }),
    $.forEach((X, Q) => {
      eY.forEach((G) => G(Q, X));
    });
  for (let X of Z) {
    if (Y.some((Q) => Q.contains(X))) continue;
    JZ.forEach((Q) => Q(X));
  }
  for (let X of Y) {
    if (!X.isConnected) continue;
    YZ.forEach((Q) => Q(X));
  }
  (Y = null), (Z = null), ($ = null), (W = null);
}
function QZ(J) {
  return AJ(YJ(J));
}
function qJ(J, Y, Z) {
  return (
    (J._x_dataStack = [Y, ...YJ(Z || J)]),
    () => {
      J._x_dataStack = J._x_dataStack.filter(($) => $ !== Y);
    }
  );
}
function YJ(J) {
  if (J._x_dataStack) return J._x_dataStack;
  if (typeof ShadowRoot === "function" && J instanceof ShadowRoot)
    return YJ(J.host);
  if (!J.parentNode) return [];
  return YJ(J.parentNode);
}
function AJ(J) {
  return new Proxy({ objects: J }, k0);
}
var k0 = {
  ownKeys({ objects: J }) {
    return Array.from(new Set(J.flatMap((Y) => Object.keys(Y))));
  },
  has({ objects: J }, Y) {
    if (Y == Symbol.unscopables) return !1;
    return J.some(
      (Z) => Object.prototype.hasOwnProperty.call(Z, Y) || Reflect.has(Z, Y),
    );
  },
  get({ objects: J }, Y, Z) {
    if (Y == "toJSON") return x0;
    return Reflect.get(J.find(($) => Reflect.has($, Y)) || {}, Y, Z);
  },
  set({ objects: J }, Y, Z, $) {
    let W =
        J.find((Q) => Object.prototype.hasOwnProperty.call(Q, Y)) ||
        J[J.length - 1],
      X = Object.getOwnPropertyDescriptor(W, Y);
    if (X?.set && X?.get) return X.set.call($, Z) || !0;
    return Reflect.set(W, Y, Z);
  },
};
function x0() {
  return Reflect.ownKeys(this).reduce((Y, Z) => {
    return (Y[Z] = Reflect.get(this, Z)), Y;
  }, {});
}
function GZ(J) {
  let Y = ($) => typeof $ === "object" && !Array.isArray($) && $ !== null,
    Z = ($, W = "") => {
      Object.entries(Object.getOwnPropertyDescriptors($)).forEach(
        ([X, { value: Q, enumerable: G }]) => {
          if (G === !1 || Q === void 0) return;
          if (typeof Q === "object" && Q !== null && Q.__v_skip) return;
          let B = W === "" ? X : `${W}.${X}`;
          if (typeof Q === "object" && Q !== null && Q._x_interceptor)
            $[X] = Q.initialize(J, B, X);
          else if (Y(Q) && Q !== $ && !(Q instanceof Element)) Z(Q, B);
        },
      );
    };
  return Z(J);
}
function BZ(J, Y = () => {}) {
  let Z = {
    initialValue: void 0,
    _x_interceptor: !0,
    initialize($, W, X) {
      return J(
        this.initialValue,
        () => S0($, W),
        (Q) => JY($, W, Q),
        W,
        X,
      );
    },
  };
  return (
    Y(Z),
    ($) => {
      if (typeof $ === "object" && $ !== null && $._x_interceptor) {
        let W = Z.initialize.bind(Z);
        Z.initialize = (X, Q, G) => {
          let B = $.initialize(X, Q, G);
          return (Z.initialValue = B), W(X, Q, G);
        };
      } else Z.initialValue = $;
      return Z;
    }
  );
}
function S0(J, Y) {
  return Y.split(".").reduce((Z, $) => Z[$], J);
}
function JY(J, Y, Z) {
  if (typeof Y === "string") Y = Y.split(".");
  if (Y.length === 1) J[Y[0]] = Z;
  else if (Y.length === 0) throw error;
  else if (J[Y[0]]) return JY(J[Y[0]], Y.slice(1), Z);
  else return (J[Y[0]] = {}), JY(J[Y[0]], Y.slice(1), Z);
}
var UZ = {};
function y(J, Y) {
  UZ[J] = Y;
}
function YY(J, Y) {
  let Z = y0(Y);
  return (
    Object.entries(UZ).forEach(([$, W]) => {
      Object.defineProperty(J, `$${$}`, {
        get() {
          return W(Y, Z);
        },
        enumerable: !1,
      });
    }),
    J
  );
}
function y0(J) {
  let [Y, Z] = qZ(J),
    $ = { interceptor: BZ, ...Y };
  return MY(J, Z), $;
}
function h0(J, Y, Z, ...$) {
  try {
    return Z(...$);
  } catch (W) {
    LJ(W, J, Y);
  }
}
function LJ(J, Y, Z = void 0) {
  (J = Object.assign(J ?? { message: "No error message given." }, {
    el: Y,
    expression: Z,
  })),
    console.warn(
      `Alpine Expression Error: ${J.message}

${
  Z
    ? 'Expression: "' +
      Z +
      `"

`
    : ""
}`,
      Y,
    ),
    setTimeout(() => {
      throw J;
    }, 0);
}
var NJ = !0;
function _Z(J) {
  let Y = NJ;
  NJ = !1;
  let Z = J();
  return (NJ = Y), Z;
}
function r(J, Y, Z = {}) {
  let $;
  return I(J, Y)((W) => ($ = W), Z), $;
}
function I(...J) {
  return zZ(...J);
}
var zZ = KZ;
function f0(J) {
  zZ = J;
}
function KZ(J, Y) {
  let Z = {};
  YY(Z, J);
  let $ = [Z, ...YJ(J)],
    W = typeof Y === "function" ? v0($, Y) : d0($, Y, J);
  return h0.bind(null, J, Y, W);
}
function v0(J, Y) {
  return (Z = () => {}, { scope: $ = {}, params: W = [] } = {}) => {
    let X = Y.apply(AJ([$, ...J]), W);
    wJ(Z, X);
  };
}
var iJ = {};
function u0(J, Y) {
  if (iJ[J]) return iJ[J];
  let Z = Object.getPrototypeOf(async function () {}).constructor,
    $ =
      /^[\n\s]*if.*\(.*\)/.test(J.trim()) || /^(let|const)\s/.test(J.trim())
        ? `(async()=>{ ${J} })()`
        : J,
    X = (() => {
      try {
        let Q = new Z(
          ["__self", "scope"],
          `with (scope) { __self.result = ${$} }; __self.finished = true; return __self.result;`,
        );
        return Object.defineProperty(Q, "name", { value: `[Alpine] ${J}` }), Q;
      } catch (Q) {
        return LJ(Q, Y, J), Promise.resolve();
      }
    })();
  return (iJ[J] = X), X;
}
function d0(J, Y, Z) {
  let $ = u0(Y, Z);
  return (W = () => {}, { scope: X = {}, params: Q = [] } = {}) => {
    ($.result = void 0), ($.finished = !1);
    let G = AJ([X, ...J]);
    if (typeof $ === "function") {
      let B = $($, G).catch((U) => LJ(U, Z, Y));
      if ($.finished) wJ(W, $.result, G, Q, Z), ($.result = void 0);
      else
        B.then((U) => {
          wJ(W, U, G, Q, Z);
        })
          .catch((U) => LJ(U, Z, Y))
          .finally(() => ($.result = void 0));
    }
  };
}
function wJ(J, Y, Z, $, W) {
  if (NJ && typeof Y === "function") {
    let X = Y.apply(Z, $);
    if (X instanceof Promise)
      X.then((Q) => wJ(J, Q, Z, $)).catch((Q) => LJ(Q, W, Y));
    else J(X);
  } else if (typeof Y === "object" && Y instanceof Promise) Y.then((X) => J(X));
  else J(Y);
}
var RY = "x-";
function XJ(J = "") {
  return RY + J;
}
function g0(J) {
  RY = J;
}
var bJ = {};
function O(J, Y) {
  return (
    (bJ[J] = Y),
    {
      before(Z) {
        if (!bJ[Z]) {
          console.warn(
            String.raw`Cannot find directive \`${Z}\`. \`${J}\` will use the default order of execution`,
          );
          return;
        }
        let $ = n.indexOf(Z);
        n.splice($ >= 0 ? $ : n.indexOf("DEFAULT"), 0, J);
      },
    }
  );
}
function m0(J) {
  return Object.keys(bJ).includes(J);
}
function HY(J, Y, Z) {
  if (((Y = Array.from(Y)), J._x_virtualDirectives)) {
    let X = Object.entries(J._x_virtualDirectives).map(([G, B]) => ({
        name: G,
        value: B,
      })),
      Q = MZ(X);
    (X = X.map((G) => {
      if (Q.find((B) => B.name === G.name))
        return { name: `x-bind:${G.name}`, value: `"${G.value}"` };
      return G;
    })),
      (Y = Y.concat(X));
  }
  let $ = {};
  return Y.map(FZ((X, Q) => ($[X] = Q)))
    .filter(HZ)
    .map(i0($, Z))
    .sort(o0)
    .map((X) => {
      return p0(J, X);
    });
}
function MZ(J) {
  return Array.from(J)
    .map(FZ())
    .filter((Y) => !HZ(Y));
}
var ZY = !1,
  zJ = new Map(),
  LZ = Symbol();
function c0(J) {
  ZY = !0;
  let Y = Symbol();
  (LZ = Y), zJ.set(Y, []);
  let Z = () => {
      while (zJ.get(Y).length) zJ.get(Y).shift()();
      zJ.delete(Y);
    },
    $ = () => {
      (ZY = !1), Z();
    };
  J(Z), $();
}
function qZ(J) {
  let Y = [],
    Z = (G) => Y.push(G),
    [$, W] = E0(J);
  return (
    Y.push(W),
    [
      {
        Alpine: PJ,
        effect: $,
        cleanup: Z,
        evaluateLater: I.bind(I, J),
        evaluate: r.bind(r, J),
      },
      () => Y.forEach((G) => G()),
    ]
  );
}
function p0(J, Y) {
  let Z = () => {},
    $ = bJ[Y.type] || Z,
    [W, X] = qZ(J);
  $Z(J, Y.original, X);
  let Q = () => {
    if (J._x_ignore || J._x_ignoreSelf) return;
    $.inline && $.inline(J, Y, W),
      ($ = $.bind($, J, Y, W)),
      ZY ? zJ.get(LZ).push($) : $();
  };
  return (Q.runCleanups = X), Q;
}
var AZ =
    (J, Y) =>
    ({ name: Z, value: $ }) => {
      if (Z.startsWith(J)) Z = Z.replace(J, Y);
      return { name: Z, value: $ };
    },
  PZ = (J) => J;
function FZ(J = () => {}) {
  return ({ name: Y, value: Z }) => {
    let { name: $, value: W } = RZ.reduce(
      (X, Q) => {
        return Q(X);
      },
      { name: Y, value: Z },
    );
    if ($ !== Y) J($, Y);
    return { name: $, value: W };
  };
}
var RZ = [];
function CY(J) {
  RZ.push(J);
}
function HZ({ name: J }) {
  return CZ().test(J);
}
var CZ = () => new RegExp(`^${RY}([^:^.]+)\\b`);
function i0(J, Y) {
  return ({ name: Z, value: $ }) => {
    let W = Z.match(CZ()),
      X = Z.match(/:([a-zA-Z0-9\-_:]+)/),
      Q = Z.match(/\.[^.\]]+(?=[^\]]*$)/g) || [],
      G = Y || J[Z] || Z;
    return {
      type: W ? W[1] : null,
      value: X ? X[1] : null,
      modifiers: Q.map((B) => B.replace(".", "")),
      expression: $,
      original: G,
    };
  };
}
var $Y = "DEFAULT",
  n = [
    "ignore",
    "ref",
    "data",
    "id",
    "anchor",
    "bind",
    "init",
    "for",
    "model",
    "modelable",
    "transition",
    "show",
    "if",
    $Y,
    "teleport",
  ];
function o0(J, Y) {
  let Z = n.indexOf(J.type) === -1 ? $Y : J.type,
    $ = n.indexOf(Y.type) === -1 ? $Y : Y.type;
  return n.indexOf(Z) - n.indexOf($);
}
function KJ(J, Y, Z = {}) {
  J.dispatchEvent(
    new CustomEvent(Y, {
      detail: Z,
      bubbles: !0,
      composed: !0,
      cancelable: !0,
    }),
  );
}
function t(J, Y) {
  if (typeof ShadowRoot === "function" && J instanceof ShadowRoot) {
    Array.from(J.children).forEach((W) => t(W, Y));
    return;
  }
  let Z = !1;
  if ((Y(J, () => (Z = !0)), Z)) return;
  let $ = J.firstElementChild;
  while ($) t($, Y, !1), ($ = $.nextElementSibling);
}
function k(J, ...Y) {
  console.warn(`Alpine Warning: ${J}`, ...Y);
}
var vY = !1;
function n0() {
  if (vY)
    k(
      "Alpine has already been initialized on this page. Calling Alpine.start() more than once can cause problems.",
    );
  if (((vY = !0), !document.body))
    k(
      "Unable to initialize. Trying to load Alpine before `<body>` is available. Did you forget to add `defer` in Alpine's `<script>` tag?",
    );
  KJ(document, "alpine:init"),
    KJ(document, "alpine:initializing"),
    AY(),
    N0((Y) => u(Y, t)),
    MY((Y) => GJ(Y)),
    ZZ((Y, Z) => {
      HY(Y, Z).forEach(($) => $());
    });
  let J = (Y) => !xJ(Y.parentElement, !0);
  Array.from(document.querySelectorAll(TZ().join(",")))
    .filter(J)
    .forEach((Y) => {
      u(Y);
    }),
    KJ(document, "alpine:initialized"),
    setTimeout(() => {
      l0();
    });
}
var jY = [],
  jZ = [];
function OZ() {
  return jY.map((J) => J());
}
function TZ() {
  return jY.concat(jZ).map((J) => J());
}
function VZ(J) {
  jY.push(J);
}
function EZ(J) {
  jZ.push(J);
}
function xJ(J, Y = !1) {
  return QJ(J, (Z) => {
    if ((Y ? TZ() : OZ()).some((W) => Z.matches(W))) return !0;
  });
}
function QJ(J, Y) {
  if (!J) return;
  if (Y(J)) return J;
  if (J._x_teleportBack) J = J._x_teleportBack;
  if (!J.parentElement) return;
  return QJ(J.parentElement, Y);
}
function s0(J) {
  return OZ().some((Y) => J.matches(Y));
}
var NZ = [];
function r0(J) {
  NZ.push(J);
}
var a0 = 1;
function u(J, Y = t, Z = () => {}) {
  if (QJ(J, ($) => $._x_ignore)) return;
  c0(() => {
    Y(J, ($, W) => {
      if ($._x_marker) return;
      if (
        (Z($, W),
        NZ.forEach((X) => X($, W)),
        HY($, $.attributes).forEach((X) => X()),
        !$._x_ignore)
      )
        $._x_marker = a0++;
      $._x_ignore && W();
    });
  });
}
function GJ(J, Y = t) {
  Y(J, (Z) => {
    I0(Z), WZ(Z), delete Z._x_marker;
  });
}
function l0() {
  [
    ["ui", "dialog", ["[x-dialog], [x-popover]"]],
    ["anchor", "anchor", ["[x-anchor]"]],
    ["sort", "sort", ["[x-sort]"]],
  ].forEach(([Y, Z, $]) => {
    if (m0(Z)) return;
    $.some((W) => {
      if (document.querySelector(W))
        return k(`found "${W}", but missing ${Y} plugin`), !0;
    });
  });
}
var WY = [],
  OY = !1;
function TY(J = () => {}) {
  return (
    queueMicrotask(() => {
      OY ||
        setTimeout(() => {
          XY();
        });
    }),
    new Promise((Y) => {
      WY.push(() => {
        J(), Y();
      });
    })
  );
}
function XY() {
  OY = !1;
  while (WY.length) WY.shift()();
}
function t0() {
  OY = !0;
}
function VY(J, Y) {
  if (Array.isArray(Y)) return uY(J, Y.join(" "));
  else if (typeof Y === "object" && Y !== null) return e0(J, Y);
  else if (typeof Y === "function") return VY(J, Y());
  return uY(J, Y);
}
function uY(J, Y) {
  let Z = (X) => X.split(" ").filter(Boolean),
    $ = (X) =>
      X.split(" ")
        .filter((Q) => !J.classList.contains(Q))
        .filter(Boolean),
    W = (X) => {
      return (
        J.classList.add(...X),
        () => {
          J.classList.remove(...X);
        }
      );
    };
  return (Y = Y === !0 ? (Y = "") : Y || ""), W($(Y));
}
function e0(J, Y) {
  let Z = (G) => G.split(" ").filter(Boolean),
    $ = Object.entries(Y)
      .flatMap(([G, B]) => (B ? Z(G) : !1))
      .filter(Boolean),
    W = Object.entries(Y)
      .flatMap(([G, B]) => (!B ? Z(G) : !1))
      .filter(Boolean),
    X = [],
    Q = [];
  return (
    W.forEach((G) => {
      if (J.classList.contains(G)) J.classList.remove(G), Q.push(G);
    }),
    $.forEach((G) => {
      if (!J.classList.contains(G)) J.classList.add(G), X.push(G);
    }),
    () => {
      Q.forEach((G) => J.classList.add(G)),
        X.forEach((G) => J.classList.remove(G));
    }
  );
}
function SJ(J, Y) {
  if (typeof Y === "object" && Y !== null) return J1(J, Y);
  return Y1(J, Y);
}
function J1(J, Y) {
  let Z = {};
  return (
    Object.entries(Y).forEach(([$, W]) => {
      if (((Z[$] = J.style[$]), !$.startsWith("--"))) $ = Z1($);
      J.style.setProperty($, W);
    }),
    setTimeout(() => {
      if (J.style.length === 0) J.removeAttribute("style");
    }),
    () => {
      SJ(J, Z);
    }
  );
}
function Y1(J, Y) {
  let Z = J.getAttribute("style", Y);
  return (
    J.setAttribute("style", Y),
    () => {
      J.setAttribute("style", Z || "");
    }
  );
}
function Z1(J) {
  return J.replace(/([a-z])([A-Z])/g, "$1-$2").toLowerCase();
}
function QY(J, Y = () => {}) {
  let Z = !1;
  return function () {
    if (!Z) (Z = !0), J.apply(this, arguments);
    else Y.apply(this, arguments);
  };
}
O(
  "transition",
  (J, { value: Y, modifiers: Z, expression: $ }, { evaluate: W }) => {
    if (typeof $ === "function") $ = W($);
    if ($ === !1) return;
    if (!$ || typeof $ === "boolean") W1(J, Z, Y);
    else $1(J, $, Y);
  },
);
function $1(J, Y, Z) {
  IZ(J, VY, ""),
    {
      enter: (W) => {
        J._x_transition.enter.during = W;
      },
      "enter-start": (W) => {
        J._x_transition.enter.start = W;
      },
      "enter-end": (W) => {
        J._x_transition.enter.end = W;
      },
      leave: (W) => {
        J._x_transition.leave.during = W;
      },
      "leave-start": (W) => {
        J._x_transition.leave.start = W;
      },
      "leave-end": (W) => {
        J._x_transition.leave.end = W;
      },
    }[Z](Y);
}
function W1(J, Y, Z) {
  IZ(J, SJ);
  let $ = !Y.includes("in") && !Y.includes("out") && !Z,
    W = $ || Y.includes("in") || ["enter"].includes(Z),
    X = $ || Y.includes("out") || ["leave"].includes(Z);
  if (Y.includes("in") && !$) Y = Y.filter((L, P) => P < Y.indexOf("out"));
  if (Y.includes("out") && !$) Y = Y.filter((L, P) => P > Y.indexOf("out"));
  let Q = !Y.includes("opacity") && !Y.includes("scale"),
    G = Q || Y.includes("opacity"),
    B = Q || Y.includes("scale"),
    U = G ? 0 : 1,
    _ = B ? UJ(Y, "scale", 95) / 100 : 1,
    z = UJ(Y, "delay", 0) / 1000,
    K = UJ(Y, "origin", "center"),
    q = "opacity, transform",
    j = UJ(Y, "duration", 150) / 1000,
    T = UJ(Y, "duration", 75) / 1000,
    M = "cubic-bezier(0.4, 0.0, 0.2, 1)";
  if (W)
    (J._x_transition.enter.during = {
      transformOrigin: K,
      transitionDelay: `${z}s`,
      transitionProperty: q,
      transitionDuration: `${j}s`,
      transitionTimingFunction: M,
    }),
      (J._x_transition.enter.start = { opacity: U, transform: `scale(${_})` }),
      (J._x_transition.enter.end = { opacity: 1, transform: "scale(1)" });
  if (X)
    (J._x_transition.leave.during = {
      transformOrigin: K,
      transitionDelay: `${z}s`,
      transitionProperty: q,
      transitionDuration: `${T}s`,
      transitionTimingFunction: M,
    }),
      (J._x_transition.leave.start = { opacity: 1, transform: "scale(1)" }),
      (J._x_transition.leave.end = { opacity: U, transform: `scale(${_})` });
}
function IZ(J, Y, Z = {}) {
  if (!J._x_transition)
    J._x_transition = {
      enter: { during: Z, start: Z, end: Z },
      leave: { during: Z, start: Z, end: Z },
      in($ = () => {}, W = () => {}) {
        GY(
          J,
          Y,
          {
            during: this.enter.during,
            start: this.enter.start,
            end: this.enter.end,
          },
          $,
          W,
        );
      },
      out($ = () => {}, W = () => {}) {
        GY(
          J,
          Y,
          {
            during: this.leave.during,
            start: this.leave.start,
            end: this.leave.end,
          },
          $,
          W,
        );
      },
    };
}
window.Element.prototype._x_toggleAndCascadeWithTransitions = function (
  J,
  Y,
  Z,
  $,
) {
  let W =
      document.visibilityState === "visible"
        ? requestAnimationFrame
        : setTimeout,
    X = () => W(Z);
  if (Y) {
    if (J._x_transition && (J._x_transition.enter || J._x_transition.leave))
      J._x_transition.enter &&
      (Object.entries(J._x_transition.enter.during).length ||
        Object.entries(J._x_transition.enter.start).length ||
        Object.entries(J._x_transition.enter.end).length)
        ? J._x_transition.in(Z)
        : X();
    else J._x_transition ? J._x_transition.in(Z) : X();
    return;
  }
  (J._x_hidePromise = J._x_transition
    ? new Promise((Q, G) => {
        J._x_transition.out(
          () => {},
          () => Q($),
        ),
          J._x_transitioning &&
            J._x_transitioning.beforeCancel(() =>
              G({ isFromCancelledTransition: !0 }),
            );
      })
    : Promise.resolve($)),
    queueMicrotask(() => {
      let Q = DZ(J);
      if (Q) {
        if (!Q._x_hideChildren) Q._x_hideChildren = [];
        Q._x_hideChildren.push(J);
      } else
        W(() => {
          let G = (B) => {
            let U = Promise.all([
              B._x_hidePromise,
              ...(B._x_hideChildren || []).map(G),
            ]).then(([_]) => _?.());
            return delete B._x_hidePromise, delete B._x_hideChildren, U;
          };
          G(J).catch((B) => {
            if (!B.isFromCancelledTransition) throw B;
          });
        });
    });
};
function DZ(J) {
  let Y = J.parentNode;
  if (!Y) return;
  return Y._x_hidePromise ? Y : DZ(Y);
}
function GY(
  J,
  Y,
  { during: Z, start: $, end: W } = {},
  X = () => {},
  Q = () => {},
) {
  if (J._x_transitioning) J._x_transitioning.cancel();
  if (
    Object.keys(Z).length === 0 &&
    Object.keys($).length === 0 &&
    Object.keys(W).length === 0
  ) {
    X(), Q();
    return;
  }
  let G, B, U;
  X1(J, {
    start() {
      G = Y(J, $);
    },
    during() {
      B = Y(J, Z);
    },
    before: X,
    end() {
      G(), (U = Y(J, W));
    },
    after: Q,
    cleanup() {
      B(), U();
    },
  });
}
function X1(J, Y) {
  let Z,
    $,
    W,
    X = QY(() => {
      C(() => {
        if (((Z = !0), !$)) Y.before();
        if (!W) Y.end(), XY();
        if ((Y.after(), J.isConnected)) Y.cleanup();
        delete J._x_transitioning;
      });
    });
  (J._x_transitioning = {
    beforeCancels: [],
    beforeCancel(Q) {
      this.beforeCancels.push(Q);
    },
    cancel: QY(function () {
      while (this.beforeCancels.length) this.beforeCancels.shift()();
      X();
    }),
    finish: X,
  }),
    C(() => {
      Y.start(), Y.during();
    }),
    t0(),
    requestAnimationFrame(() => {
      if (Z) return;
      let Q =
          Number(
            getComputedStyle(J)
              .transitionDuration.replace(/,.*/, "")
              .replace("s", ""),
          ) * 1000,
        G =
          Number(
            getComputedStyle(J)
              .transitionDelay.replace(/,.*/, "")
              .replace("s", ""),
          ) * 1000;
      if (Q === 0)
        Q =
          Number(getComputedStyle(J).animationDuration.replace("s", "")) * 1000;
      C(() => {
        Y.before();
      }),
        ($ = !0),
        requestAnimationFrame(() => {
          if (Z) return;
          C(() => {
            Y.end();
          }),
            XY(),
            setTimeout(J._x_transitioning.finish, Q + G),
            (W = !0);
        });
    });
}
function UJ(J, Y, Z) {
  if (J.indexOf(Y) === -1) return Z;
  let $ = J[J.indexOf(Y) + 1];
  if (!$) return Z;
  if (Y === "scale") {
    if (isNaN($)) return Z;
  }
  if (Y === "duration" || Y === "delay") {
    let W = $.match(/([0-9]+)ms/);
    if (W) return W[1];
  }
  if (Y === "origin") {
    if (
      ["top", "right", "left", "center", "bottom"].includes(J[J.indexOf(Y) + 2])
    )
      return [$, J[J.indexOf(Y) + 2]].join(" ");
  }
  return $;
}
var m = !1;
function p(J, Y = () => {}) {
  return (...Z) => (m ? Y(...Z) : J(...Z));
}
function Q1(J) {
  return (...Y) => m && J(...Y);
}
var wZ = [];
function yJ(J) {
  wZ.push(J);
}
function G1(J, Y) {
  wZ.forEach((Z) => Z(J, Y)),
    (m = !0),
    bZ(() => {
      u(Y, (Z, $) => {
        $(Z, () => {});
      });
    }),
    (m = !1);
}
var BY = !1;
function B1(J, Y) {
  if (!Y._x_dataStack) Y._x_dataStack = J._x_dataStack;
  (m = !0),
    (BY = !0),
    bZ(() => {
      U1(Y);
    }),
    (m = !1),
    (BY = !1);
}
function U1(J) {
  let Y = !1;
  u(J, ($, W) => {
    t($, (X, Q) => {
      if (Y && s0(X)) return Q();
      (Y = !0), W(X, Q);
    });
  });
}
function bZ(J) {
  let Y = e;
  fY((Z, $) => {
    let W = Y(Z);
    return WJ(W), () => {};
  }),
    J(),
    fY(Y);
}
function kZ(J, Y, Z, $ = []) {
  if (!J._x_bindings) J._x_bindings = $J({});
  switch (((J._x_bindings[Y] = Z), (Y = $.includes("camel") ? P1(Y) : Y), Y)) {
    case "value":
      _1(J, Z);
      break;
    case "style":
      K1(J, Z);
      break;
    case "class":
      z1(J, Z);
      break;
    case "selected":
    case "checked":
      M1(J, Y, Z);
      break;
    default:
      xZ(J, Y, Z);
      break;
  }
}
function _1(J, Y) {
  if (hZ(J)) {
    if (J.attributes.value === void 0) J.value = Y;
    if (window.fromModel)
      if (typeof Y === "boolean") J.checked = IJ(J.value) === Y;
      else J.checked = dY(J.value, Y);
  } else if (EY(J))
    if (Number.isInteger(Y)) J.value = Y;
    else if (
      !Array.isArray(Y) &&
      typeof Y !== "boolean" &&
      ![null, void 0].includes(Y)
    )
      J.value = String(Y);
    else if (Array.isArray(Y)) J.checked = Y.some((Z) => dY(Z, J.value));
    else J.checked = !!Y;
  else if (J.tagName === "SELECT") A1(J, Y);
  else {
    if (J.value === Y) return;
    J.value = Y === void 0 ? "" : Y;
  }
}
function z1(J, Y) {
  if (J._x_undoAddedClasses) J._x_undoAddedClasses();
  J._x_undoAddedClasses = VY(J, Y);
}
function K1(J, Y) {
  if (J._x_undoAddedStyles) J._x_undoAddedStyles();
  J._x_undoAddedStyles = SJ(J, Y);
}
function M1(J, Y, Z) {
  xZ(J, Y, Z), q1(J, Y, Z);
}
function xZ(J, Y, Z) {
  if ([null, void 0, !1].includes(Z) && R1(Y)) J.removeAttribute(Y);
  else {
    if (SZ(Y)) Z = Y;
    L1(J, Y, Z);
  }
}
function L1(J, Y, Z) {
  if (J.getAttribute(Y) != Z) J.setAttribute(Y, Z);
}
function q1(J, Y, Z) {
  if (J[Y] !== Z) J[Y] = Z;
}
function A1(J, Y) {
  let Z = [].concat(Y).map(($) => {
    return $ + "";
  });
  Array.from(J.options).forEach(($) => {
    $.selected = Z.includes($.value);
  });
}
function P1(J) {
  return J.toLowerCase().replace(/-(\w)/g, (Y, Z) => Z.toUpperCase());
}
function dY(J, Y) {
  return J == Y;
}
function IJ(J) {
  if ([1, "1", "true", "on", "yes", !0].includes(J)) return !0;
  if ([0, "0", "false", "off", "no", !1].includes(J)) return !1;
  return J ? Boolean(J) : null;
}
var F1 = new Set([
  "allowfullscreen",
  "async",
  "autofocus",
  "autoplay",
  "checked",
  "controls",
  "default",
  "defer",
  "disabled",
  "formnovalidate",
  "inert",
  "ismap",
  "itemscope",
  "loop",
  "multiple",
  "muted",
  "nomodule",
  "novalidate",
  "open",
  "playsinline",
  "readonly",
  "required",
  "reversed",
  "selected",
  "shadowrootclonable",
  "shadowrootdelegatesfocus",
  "shadowrootserializable",
]);
function SZ(J) {
  return F1.has(J);
}
function R1(J) {
  return ![
    "aria-pressed",
    "aria-checked",
    "aria-expanded",
    "aria-selected",
  ].includes(J);
}
function H1(J, Y, Z) {
  if (J._x_bindings && J._x_bindings[Y] !== void 0) return J._x_bindings[Y];
  return yZ(J, Y, Z);
}
function C1(J, Y, Z, $ = !0) {
  if (J._x_bindings && J._x_bindings[Y] !== void 0) return J._x_bindings[Y];
  if (J._x_inlineBindings && J._x_inlineBindings[Y] !== void 0) {
    let W = J._x_inlineBindings[Y];
    return (
      (W.extract = $),
      _Z(() => {
        return r(J, W.expression);
      })
    );
  }
  return yZ(J, Y, Z);
}
function yZ(J, Y, Z) {
  let $ = J.getAttribute(Y);
  if ($ === null) return typeof Z === "function" ? Z() : Z;
  if ($ === "") return !0;
  if (SZ(Y)) return !![Y, "true"].includes($);
  return $;
}
function EY(J) {
  return (
    J.type === "checkbox" ||
    J.localName === "ui-checkbox" ||
    J.localName === "ui-switch"
  );
}
function hZ(J) {
  return J.type === "radio" || J.localName === "ui-radio";
}
function fZ(J, Y) {
  var Z;
  return function () {
    var $ = this,
      W = arguments,
      X = function () {
        (Z = null), J.apply($, W);
      };
    clearTimeout(Z), (Z = setTimeout(X, Y));
  };
}
function vZ(J, Y) {
  let Z;
  return function () {
    let $ = this,
      W = arguments;
    if (!Z) J.apply($, W), (Z = !0), setTimeout(() => (Z = !1), Y);
  };
}
function uZ({ get: J, set: Y }, { get: Z, set: $ }) {
  let W = !0,
    X,
    Q,
    G = e(() => {
      let B = J(),
        U = Z();
      if (W) $(oJ(B)), (W = !1);
      else {
        let _ = JSON.stringify(B),
          z = JSON.stringify(U);
        if (_ !== X) $(oJ(B));
        else if (_ !== z) Y(oJ(U));
      }
      (X = JSON.stringify(J())), (Q = JSON.stringify(Z()));
    });
  return () => {
    WJ(G);
  };
}
function oJ(J) {
  return typeof J === "object" ? JSON.parse(JSON.stringify(J)) : J;
}
function j1(J) {
  (Array.isArray(J) ? J : [J]).forEach((Z) => Z(PJ));
}
var o = {},
  gY = !1;
function O1(J, Y) {
  if (!gY) (o = $J(o)), (gY = !0);
  if (Y === void 0) return o[J];
  if (
    ((o[J] = Y),
    GZ(o[J]),
    typeof Y === "object" &&
      Y !== null &&
      Y.hasOwnProperty("init") &&
      typeof Y.init === "function")
  )
    o[J].init();
}
function T1() {
  return o;
}
var dZ = {};
function V1(J, Y) {
  let Z = typeof Y !== "function" ? () => Y : Y;
  if (J instanceof Element) return gZ(J, Z());
  else dZ[J] = Z;
  return () => {};
}
function E1(J) {
  return (
    Object.entries(dZ).forEach(([Y, Z]) => {
      Object.defineProperty(J, Y, {
        get() {
          return (...$) => {
            return Z(...$);
          };
        },
      });
    }),
    J
  );
}
function gZ(J, Y, Z) {
  let $ = [];
  while ($.length) $.pop()();
  let W = Object.entries(Y).map(([Q, G]) => ({ name: Q, value: G })),
    X = MZ(W);
  return (
    (W = W.map((Q) => {
      if (X.find((G) => G.name === Q.name))
        return { name: `x-bind:${Q.name}`, value: `"${Q.value}"` };
      return Q;
    })),
    HY(J, W, Z).map((Q) => {
      $.push(Q.runCleanups), Q();
    }),
    () => {
      while ($.length) $.pop()();
    }
  );
}
var mZ = {};
function N1(J, Y) {
  mZ[J] = Y;
}
function I1(J, Y) {
  return (
    Object.entries(mZ).forEach(([Z, $]) => {
      Object.defineProperty(J, Z, {
        get() {
          return (...W) => {
            return $.bind(Y)(...W);
          };
        },
        enumerable: !1,
      });
    }),
    J
  );
}
var D1 = {
    get reactive() {
      return $J;
    },
    get release() {
      return WJ;
    },
    get effect() {
      return e;
    },
    get raw() {
      return lY;
    },
    version: "3.14.9",
    flushAndStopDeferringMutations: b0,
    dontAutoEvaluateFunctions: _Z,
    disableEffectScheduling: T0,
    startObservingMutations: AY,
    stopObservingMutations: XZ,
    setReactivityEngine: V0,
    onAttributeRemoved: $Z,
    onAttributesAdded: ZZ,
    closestDataStack: YJ,
    skipDuringClone: p,
    onlyDuringClone: Q1,
    addRootSelector: VZ,
    addInitSelector: EZ,
    interceptClone: yJ,
    addScopeToNode: qJ,
    deferMutations: w0,
    mapAttributes: CY,
    evaluateLater: I,
    interceptInit: r0,
    setEvaluator: f0,
    mergeProxies: AJ,
    extractProp: C1,
    findClosest: QJ,
    onElRemoved: MY,
    closestRoot: xJ,
    destroyTree: GJ,
    interceptor: BZ,
    transition: GY,
    setStyles: SJ,
    mutateDom: C,
    directive: O,
    entangle: uZ,
    throttle: vZ,
    debounce: fZ,
    evaluate: r,
    initTree: u,
    nextTick: TY,
    prefixed: XJ,
    prefix: g0,
    plugin: j1,
    magic: y,
    store: O1,
    start: n0,
    clone: B1,
    cloneNode: G1,
    bound: H1,
    $data: QZ,
    watch: tY,
    walk: t,
    data: N1,
    bind: V1,
  },
  PJ = D1;
function cZ(J, Y) {
  let Z = Object.create(null),
    $ = J.split(",");
  for (let W = 0; W < $.length; W++) Z[$[W]] = !0;
  return Y ? (W) => !!Z[W.toLowerCase()] : (W) => !!Z[W];
}
var w1 =
    "itemscope,allowfullscreen,formnovalidate,ismap,nomodule,novalidate,readonly",
  x6 = cZ(
    w1 +
      ",async,autofocus,autoplay,controls,default,defer,disabled,hidden,loop,open,required,reversed,scoped,seamless,checked,muted,multiple,selected",
  ),
  b1 = Object.freeze({}),
  S6 = Object.freeze([]),
  k1 = Object.prototype.hasOwnProperty,
  hJ = (J, Y) => k1.call(J, Y),
  a = Array.isArray,
  MJ = (J) => pZ(J) === "[object Map]",
  x1 = (J) => typeof J === "string",
  NY = (J) => typeof J === "symbol",
  fJ = (J) => J !== null && typeof J === "object",
  S1 = Object.prototype.toString,
  pZ = (J) => S1.call(J),
  iZ = (J) => {
    return pZ(J).slice(8, -1);
  },
  IY = (J) =>
    x1(J) && J !== "NaN" && J[0] !== "-" && "" + parseInt(J, 10) === J,
  vJ = (J) => {
    let Y = Object.create(null);
    return (Z) => {
      return Y[Z] || (Y[Z] = J(Z));
    };
  },
  y1 = /-(\w)/g,
  y6 = vJ((J) => {
    return J.replace(y1, (Y, Z) => (Z ? Z.toUpperCase() : ""));
  }),
  h1 = /\B([A-Z])/g,
  h6 = vJ((J) => J.replace(h1, "-$1").toLowerCase()),
  oZ = vJ((J) => J.charAt(0).toUpperCase() + J.slice(1)),
  f6 = vJ((J) => (J ? `on${oZ(J)}` : "")),
  nZ = (J, Y) => J !== Y && (J === J || Y === Y),
  UY = new WeakMap(),
  _J = [],
  h,
  l = Symbol("iterate"),
  _Y = Symbol("Map key iterate");
function f1(J) {
  return J && J._isEffect === !0;
}
function v1(J, Y = b1) {
  if (f1(J)) J = J.raw;
  let Z = g1(J, Y);
  if (!Y.lazy) Z();
  return Z;
}
function u1(J) {
  if (J.active) {
    if ((sZ(J), J.options.onStop)) J.options.onStop();
    J.active = !1;
  }
}
var d1 = 0;
function g1(J, Y) {
  let Z = function $() {
    if (!Z.active) return J();
    if (!_J.includes(Z)) {
      sZ(Z);
      try {
        return c1(), _J.push(Z), (h = Z), J();
      } finally {
        _J.pop(), rZ(), (h = _J[_J.length - 1]);
      }
    }
  };
  return (
    (Z.id = d1++),
    (Z.allowRecurse = !!Y.allowRecurse),
    (Z._isEffect = !0),
    (Z.active = !0),
    (Z.raw = J),
    (Z.deps = []),
    (Z.options = Y),
    Z
  );
}
function sZ(J) {
  let { deps: Y } = J;
  if (Y.length) {
    for (let Z = 0; Z < Y.length; Z++) Y[Z].delete(J);
    Y.length = 0;
  }
}
var ZJ = !0,
  DY = [];
function m1() {
  DY.push(ZJ), (ZJ = !1);
}
function c1() {
  DY.push(ZJ), (ZJ = !0);
}
function rZ() {
  let J = DY.pop();
  ZJ = J === void 0 ? !0 : J;
}
function S(J, Y, Z) {
  if (!ZJ || h === void 0) return;
  let $ = UY.get(J);
  if (!$) UY.set(J, ($ = new Map()));
  let W = $.get(Z);
  if (!W) $.set(Z, (W = new Set()));
  if (!W.has(h)) {
    if ((W.add(h), h.deps.push(W), h.options.onTrack))
      h.options.onTrack({ effect: h, target: J, type: Y, key: Z });
  }
}
function c(J, Y, Z, $, W, X) {
  let Q = UY.get(J);
  if (!Q) return;
  let G = new Set(),
    B = (_) => {
      if (_)
        _.forEach((z) => {
          if (z !== h || z.allowRecurse) G.add(z);
        });
    };
  if (Y === "clear") Q.forEach(B);
  else if (Z === "length" && a(J))
    Q.forEach((_, z) => {
      if (z === "length" || z >= $) B(_);
    });
  else {
    if (Z !== void 0) B(Q.get(Z));
    switch (Y) {
      case "add":
        if (!a(J)) {
          if ((B(Q.get(l)), MJ(J))) B(Q.get(_Y));
        } else if (IY(Z)) B(Q.get("length"));
        break;
      case "delete":
        if (!a(J)) {
          if ((B(Q.get(l)), MJ(J))) B(Q.get(_Y));
        }
        break;
      case "set":
        if (MJ(J)) B(Q.get(l));
        break;
    }
  }
  let U = (_) => {
    if (_.options.onTrigger)
      _.options.onTrigger({
        effect: _,
        target: J,
        key: Z,
        type: Y,
        newValue: $,
        oldValue: W,
        oldTarget: X,
      });
    if (_.options.scheduler) _.options.scheduler(_);
    else _();
  };
  G.forEach(U);
}
var p1 = cZ("__proto__,__v_isRef,__isVue"),
  aZ = new Set(
    Object.getOwnPropertyNames(Symbol)
      .map((J) => Symbol[J])
      .filter(NY),
  ),
  i1 = lZ(),
  o1 = lZ(!0),
  mY = n1();
function n1() {
  let J = {};
  return (
    ["includes", "indexOf", "lastIndexOf"].forEach((Y) => {
      J[Y] = function (...Z) {
        let $ = R(this);
        for (let X = 0, Q = this.length; X < Q; X++) S($, "get", X + "");
        let W = $[Y](...Z);
        if (W === -1 || W === !1) return $[Y](...Z.map(R));
        else return W;
      };
    }),
    ["push", "pop", "shift", "unshift", "splice"].forEach((Y) => {
      J[Y] = function (...Z) {
        m1();
        let $ = R(this)[Y].apply(this, Z);
        return rZ(), $;
      };
    }),
    J
  );
}
function lZ(J = !1, Y = !1) {
  return function Z($, W, X) {
    if (W === "__v_isReactive") return !J;
    else if (W === "__v_isReadonly") return J;
    else if (W === "__v_raw" && X === (J ? (Y ? U6 : Y0) : Y ? B6 : J0).get($))
      return $;
    let Q = a($);
    if (!J && Q && hJ(mY, W)) return Reflect.get(mY, W, X);
    let G = Reflect.get($, W, X);
    if (NY(W) ? aZ.has(W) : p1(W)) return G;
    if (!J) S($, "get", W);
    if (Y) return G;
    if (zY(G)) return !Q || !IY(W) ? G.value : G;
    if (fJ(G)) return J ? Z0(G) : xY(G);
    return G;
  };
}
var s1 = r1();
function r1(J = !1) {
  return function Y(Z, $, W, X) {
    let Q = Z[$];
    if (!J) {
      if (((W = R(W)), (Q = R(Q)), !a(Z) && zY(Q) && !zY(W)))
        return (Q.value = W), !0;
    }
    let G = a(Z) && IY($) ? Number($) < Z.length : hJ(Z, $),
      B = Reflect.set(Z, $, W, X);
    if (Z === R(X)) {
      if (!G) c(Z, "add", $, W);
      else if (nZ(W, Q)) c(Z, "set", $, W, Q);
    }
    return B;
  };
}
function a1(J, Y) {
  let Z = hJ(J, Y),
    $ = J[Y],
    W = Reflect.deleteProperty(J, Y);
  if (W && Z) c(J, "delete", Y, void 0, $);
  return W;
}
function l1(J, Y) {
  let Z = Reflect.has(J, Y);
  if (!NY(Y) || !aZ.has(Y)) S(J, "has", Y);
  return Z;
}
function t1(J) {
  return S(J, "iterate", a(J) ? "length" : l), Reflect.ownKeys(J);
}
var e1 = { get: i1, set: s1, deleteProperty: a1, has: l1, ownKeys: t1 },
  J6 = {
    get: o1,
    set(J, Y) {
      return (
        console.warn(
          `Set operation on key "${String(Y)}" failed: target is readonly.`,
          J,
        ),
        !0
      );
    },
    deleteProperty(J, Y) {
      return (
        console.warn(
          `Delete operation on key "${String(Y)}" failed: target is readonly.`,
          J,
        ),
        !0
      );
    },
  },
  wY = (J) => (fJ(J) ? xY(J) : J),
  bY = (J) => (fJ(J) ? Z0(J) : J),
  kY = (J) => J,
  uJ = (J) => Reflect.getPrototypeOf(J);
function jJ(J, Y, Z = !1, $ = !1) {
  J = J.__v_raw;
  let W = R(J),
    X = R(Y);
  if (Y !== X) !Z && S(W, "get", Y);
  !Z && S(W, "get", X);
  let { has: Q } = uJ(W),
    G = $ ? kY : Z ? bY : wY;
  if (Q.call(W, Y)) return G(J.get(Y));
  else if (Q.call(W, X)) return G(J.get(X));
  else if (J !== W) J.get(Y);
}
function OJ(J, Y = !1) {
  let Z = this.__v_raw,
    $ = R(Z),
    W = R(J);
  if (J !== W) !Y && S($, "has", J);
  return !Y && S($, "has", W), J === W ? Z.has(J) : Z.has(J) || Z.has(W);
}
function TJ(J, Y = !1) {
  return (
    (J = J.__v_raw), !Y && S(R(J), "iterate", l), Reflect.get(J, "size", J)
  );
}
function cY(J) {
  J = R(J);
  let Y = R(this);
  if (!uJ(Y).has.call(Y, J)) Y.add(J), c(Y, "add", J, J);
  return this;
}
function pY(J, Y) {
  Y = R(Y);
  let Z = R(this),
    { has: $, get: W } = uJ(Z),
    X = $.call(Z, J);
  if (!X) (J = R(J)), (X = $.call(Z, J));
  else eZ(Z, $, J);
  let Q = W.call(Z, J);
  if ((Z.set(J, Y), !X)) c(Z, "add", J, Y);
  else if (nZ(Y, Q)) c(Z, "set", J, Y, Q);
  return this;
}
function iY(J) {
  let Y = R(this),
    { has: Z, get: $ } = uJ(Y),
    W = Z.call(Y, J);
  if (!W) (J = R(J)), (W = Z.call(Y, J));
  else eZ(Y, Z, J);
  let X = $ ? $.call(Y, J) : void 0,
    Q = Y.delete(J);
  if (W) c(Y, "delete", J, void 0, X);
  return Q;
}
function oY() {
  let J = R(this),
    Y = J.size !== 0,
    Z = MJ(J) ? new Map(J) : new Set(J),
    $ = J.clear();
  if (Y) c(J, "clear", void 0, void 0, Z);
  return $;
}
function VJ(J, Y) {
  return function Z($, W) {
    let X = this,
      Q = X.__v_raw,
      G = R(Q),
      B = Y ? kY : J ? bY : wY;
    return (
      !J && S(G, "iterate", l),
      Q.forEach((U, _) => {
        return $.call(W, B(U), B(_), X);
      })
    );
  };
}
function EJ(J, Y, Z) {
  return function (...$) {
    let W = this.__v_raw,
      X = R(W),
      Q = MJ(X),
      G = J === "entries" || (J === Symbol.iterator && Q),
      B = J === "keys" && Q,
      U = W[J](...$),
      _ = Z ? kY : Y ? bY : wY;
    return (
      !Y && S(X, "iterate", B ? _Y : l),
      {
        next() {
          let { value: z, done: K } = U.next();
          return K
            ? { value: z, done: K }
            : { value: G ? [_(z[0]), _(z[1])] : _(z), done: K };
        },
        [Symbol.iterator]() {
          return this;
        },
      }
    );
  };
}
function g(J) {
  return function (...Y) {
    {
      let Z = Y[0] ? `on key "${Y[0]}" ` : "";
      console.warn(
        `${oZ(J)} operation ${Z}failed: target is readonly.`,
        R(this),
      );
    }
    return J === "delete" ? !1 : this;
  };
}
function Y6() {
  let J = {
      get(X) {
        return jJ(this, X);
      },
      get size() {
        return TJ(this);
      },
      has: OJ,
      add: cY,
      set: pY,
      delete: iY,
      clear: oY,
      forEach: VJ(!1, !1),
    },
    Y = {
      get(X) {
        return jJ(this, X, !1, !0);
      },
      get size() {
        return TJ(this);
      },
      has: OJ,
      add: cY,
      set: pY,
      delete: iY,
      clear: oY,
      forEach: VJ(!1, !0),
    },
    Z = {
      get(X) {
        return jJ(this, X, !0);
      },
      get size() {
        return TJ(this, !0);
      },
      has(X) {
        return OJ.call(this, X, !0);
      },
      add: g("add"),
      set: g("set"),
      delete: g("delete"),
      clear: g("clear"),
      forEach: VJ(!0, !1),
    },
    $ = {
      get(X) {
        return jJ(this, X, !0, !0);
      },
      get size() {
        return TJ(this, !0);
      },
      has(X) {
        return OJ.call(this, X, !0);
      },
      add: g("add"),
      set: g("set"),
      delete: g("delete"),
      clear: g("clear"),
      forEach: VJ(!0, !0),
    };
  return (
    ["keys", "values", "entries", Symbol.iterator].forEach((X) => {
      (J[X] = EJ(X, !1, !1)),
        (Z[X] = EJ(X, !0, !1)),
        (Y[X] = EJ(X, !1, !0)),
        ($[X] = EJ(X, !0, !0));
    }),
    [J, Z, Y, $]
  );
}
var [Z6, $6, W6, X6] = Y6();
function tZ(J, Y) {
  let Z = Y ? (J ? X6 : W6) : J ? $6 : Z6;
  return ($, W, X) => {
    if (W === "__v_isReactive") return !J;
    else if (W === "__v_isReadonly") return J;
    else if (W === "__v_raw") return $;
    return Reflect.get(hJ(Z, W) && W in $ ? Z : $, W, X);
  };
}
var Q6 = { get: tZ(!1, !1) },
  G6 = { get: tZ(!0, !1) };
function eZ(J, Y, Z) {
  let $ = R(Z);
  if ($ !== Z && Y.call(J, $)) {
    let W = iZ(J);
    console.warn(
      `Reactive ${W} contains both the raw and reactive versions of the same object${W === "Map" ? " as keys" : ""}, which can lead to inconsistencies. Avoid differentiating between the raw and reactive versions of an object and only use the reactive version if possible.`,
    );
  }
}
var J0 = new WeakMap(),
  B6 = new WeakMap(),
  Y0 = new WeakMap(),
  U6 = new WeakMap();
function _6(J) {
  switch (J) {
    case "Object":
    case "Array":
      return 1;
    case "Map":
    case "Set":
    case "WeakMap":
    case "WeakSet":
      return 2;
    default:
      return 0;
  }
}
function z6(J) {
  return J.__v_skip || !Object.isExtensible(J) ? 0 : _6(iZ(J));
}
function xY(J) {
  if (J && J.__v_isReadonly) return J;
  return $0(J, !1, e1, Q6, J0);
}
function Z0(J) {
  return $0(J, !0, J6, G6, Y0);
}
function $0(J, Y, Z, $, W) {
  if (!fJ(J))
    return console.warn(`value cannot be made reactive: ${String(J)}`), J;
  if (J.__v_raw && !(Y && J.__v_isReactive)) return J;
  let X = W.get(J);
  if (X) return X;
  let Q = z6(J);
  if (Q === 0) return J;
  let G = new Proxy(J, Q === 2 ? $ : Z);
  return W.set(J, G), G;
}
function R(J) {
  return (J && R(J.__v_raw)) || J;
}
function zY(J) {
  return Boolean(J && J.__v_isRef === !0);
}
y("nextTick", () => TY);
y("dispatch", (J) => KJ.bind(KJ, J));
y("watch", (J, { evaluateLater: Y, cleanup: Z }) => ($, W) => {
  let X = Y($),
    G = tY(() => {
      let B;
      return X((U) => (B = U)), B;
    }, W);
  Z(G);
});
y("store", T1);
y("data", (J) => QZ(J));
y("root", (J) => xJ(J));
y("refs", (J) => {
  if (J._x_refs_proxy) return J._x_refs_proxy;
  return (J._x_refs_proxy = AJ(K6(J))), J._x_refs_proxy;
});
function K6(J) {
  let Y = [];
  return (
    QJ(J, (Z) => {
      if (Z._x_refs) Y.push(Z._x_refs);
    }),
    Y
  );
}
var nJ = {};
function W0(J) {
  if (!nJ[J]) nJ[J] = 0;
  return ++nJ[J];
}
function M6(J, Y) {
  return QJ(J, (Z) => {
    if (Z._x_ids && Z._x_ids[Y]) return !0;
  });
}
function L6(J, Y) {
  if (!J._x_ids) J._x_ids = {};
  if (!J._x_ids[Y]) J._x_ids[Y] = W0(Y);
}
y("id", (J, { cleanup: Y }) => (Z, $ = null) => {
  let W = `${Z}${$ ? `-${$}` : ""}`;
  return q6(J, W, Y, () => {
    let X = M6(J, Z),
      Q = X ? X._x_ids[Z] : W0(Z);
    return $ ? `${Z}-${Q}-${$}` : `${Z}-${Q}`;
  });
});
yJ((J, Y) => {
  if (J._x_id) Y._x_id = J._x_id;
});
function q6(J, Y, Z, $) {
  if (!J._x_id) J._x_id = {};
  if (J._x_id[Y]) return J._x_id[Y];
  let W = $();
  return (
    (J._x_id[Y] = W),
    Z(() => {
      delete J._x_id[Y];
    }),
    W
  );
}
y("el", (J) => J);
X0("Focus", "focus", "focus");
X0("Persist", "persist", "persist");
function X0(J, Y, Z) {
  y(Y, ($) =>
    k(
      `You can't use [$${Y}] without first installing the "${J}" plugin here: https://alpinejs.dev/plugins/${Z}`,
      $,
    ),
  );
}
O(
  "modelable",
  (J, { expression: Y }, { effect: Z, evaluateLater: $, cleanup: W }) => {
    let X = $(Y),
      Q = () => {
        let _;
        return X((z) => (_ = z)), _;
      },
      G = $(`${Y} = __placeholder`),
      B = (_) => G(() => {}, { scope: { __placeholder: _ } }),
      U = Q();
    B(U),
      queueMicrotask(() => {
        if (!J._x_model) return;
        J._x_removeModelListeners.default();
        let _ = J._x_model.get,
          z = J._x_model.set,
          K = uZ(
            {
              get() {
                return _();
              },
              set(q) {
                z(q);
              },
            },
            {
              get() {
                return Q();
              },
              set(q) {
                B(q);
              },
            },
          );
        W(K);
      });
  },
);
O("teleport", (J, { modifiers: Y, expression: Z }, { cleanup: $ }) => {
  if (J.tagName.toLowerCase() !== "template")
    k("x-teleport can only be used on a <template> tag", J);
  let W = nY(Z),
    X = J.content.cloneNode(!0).firstElementChild;
  if (
    ((J._x_teleport = X),
    (X._x_teleportBack = J),
    J.setAttribute("data-teleport-template", !0),
    X.setAttribute("data-teleport-target", !0),
    J._x_forwardEvents)
  )
    J._x_forwardEvents.forEach((G) => {
      X.addEventListener(G, (B) => {
        B.stopPropagation(), J.dispatchEvent(new B.constructor(B.type, B));
      });
    });
  qJ(X, {}, J);
  let Q = (G, B, U) => {
    if (U.includes("prepend")) B.parentNode.insertBefore(G, B);
    else if (U.includes("append")) B.parentNode.insertBefore(G, B.nextSibling);
    else B.appendChild(G);
  };
  C(() => {
    Q(X, W, Y),
      p(() => {
        u(X);
      })();
  }),
    (J._x_teleportPutBack = () => {
      let G = nY(Z);
      C(() => {
        Q(J._x_teleport, G, Y);
      });
    }),
    $(() =>
      C(() => {
        X.remove(), GJ(X);
      }),
    );
});
var A6 = document.createElement("div");
function nY(J) {
  let Y = p(
    () => {
      return document.querySelector(J);
    },
    () => {
      return A6;
    },
  )();
  if (!Y) k(`Cannot find x-teleport element for selector: "${J}"`);
  return Y;
}
var Q0 = () => {};
Q0.inline = (J, { modifiers: Y }, { cleanup: Z }) => {
  Y.includes("self") ? (J._x_ignoreSelf = !0) : (J._x_ignore = !0),
    Z(() => {
      Y.includes("self") ? delete J._x_ignoreSelf : delete J._x_ignore;
    });
};
O("ignore", Q0);
O(
  "effect",
  p((J, { expression: Y }, { effect: Z }) => {
    Z(I(J, Y));
  }),
);
function KY(J, Y, Z, $) {
  let W = J,
    X = (B) => $(B),
    Q = {},
    G = (B, U) => (_) => U(B, _);
  if (Z.includes("dot")) Y = P6(Y);
  if (Z.includes("camel")) Y = F6(Y);
  if (Z.includes("passive")) Q.passive = !0;
  if (Z.includes("capture")) Q.capture = !0;
  if (Z.includes("window")) W = window;
  if (Z.includes("document")) W = document;
  if (Z.includes("debounce")) {
    let B = Z[Z.indexOf("debounce") + 1] || "invalid-wait",
      U = kJ(B.split("ms")[0]) ? Number(B.split("ms")[0]) : 250;
    X = fZ(X, U);
  }
  if (Z.includes("throttle")) {
    let B = Z[Z.indexOf("throttle") + 1] || "invalid-wait",
      U = kJ(B.split("ms")[0]) ? Number(B.split("ms")[0]) : 250;
    X = vZ(X, U);
  }
  if (Z.includes("prevent"))
    X = G(X, (B, U) => {
      U.preventDefault(), B(U);
    });
  if (Z.includes("stop"))
    X = G(X, (B, U) => {
      U.stopPropagation(), B(U);
    });
  if (Z.includes("once"))
    X = G(X, (B, U) => {
      B(U), W.removeEventListener(Y, X, Q);
    });
  if (Z.includes("away") || Z.includes("outside"))
    (W = document),
      (X = G(X, (B, U) => {
        if (J.contains(U.target)) return;
        if (U.target.isConnected === !1) return;
        if (J.offsetWidth < 1 && J.offsetHeight < 1) return;
        if (J._x_isShown === !1) return;
        B(U);
      }));
  if (Z.includes("self"))
    X = G(X, (B, U) => {
      U.target === J && B(U);
    });
  if (H6(Y) || G0(Y))
    X = G(X, (B, U) => {
      if (C6(U, Z)) return;
      B(U);
    });
  return (
    W.addEventListener(Y, X, Q),
    () => {
      W.removeEventListener(Y, X, Q);
    }
  );
}
function P6(J) {
  return J.replace(/-/g, ".");
}
function F6(J) {
  return J.toLowerCase().replace(/-(\w)/g, (Y, Z) => Z.toUpperCase());
}
function kJ(J) {
  return !Array.isArray(J) && !isNaN(J);
}
function R6(J) {
  if ([" ", "_"].includes(J)) return J;
  return J.replace(/([a-z])([A-Z])/g, "$1-$2")
    .replace(/[_\s]/, "-")
    .toLowerCase();
}
function H6(J) {
  return ["keydown", "keyup"].includes(J);
}
function G0(J) {
  return ["contextmenu", "click", "mouse"].some((Y) => J.includes(Y));
}
function C6(J, Y) {
  let Z = Y.filter((X) => {
    return ![
      "window",
      "document",
      "prevent",
      "stop",
      "once",
      "capture",
      "self",
      "away",
      "outside",
      "passive",
    ].includes(X);
  });
  if (Z.includes("debounce")) {
    let X = Z.indexOf("debounce");
    Z.splice(X, kJ((Z[X + 1] || "invalid-wait").split("ms")[0]) ? 2 : 1);
  }
  if (Z.includes("throttle")) {
    let X = Z.indexOf("throttle");
    Z.splice(X, kJ((Z[X + 1] || "invalid-wait").split("ms")[0]) ? 2 : 1);
  }
  if (Z.length === 0) return !1;
  if (Z.length === 1 && sY(J.key).includes(Z[0])) return !1;
  let W = ["ctrl", "shift", "alt", "meta", "cmd", "super"].filter((X) =>
    Z.includes(X),
  );
  if (((Z = Z.filter((X) => !W.includes(X))), W.length > 0)) {
    if (
      W.filter((Q) => {
        if (Q === "cmd" || Q === "super") Q = "meta";
        return J[`${Q}Key`];
      }).length === W.length
    ) {
      if (G0(J.type)) return !1;
      if (sY(J.key).includes(Z[0])) return !1;
    }
  }
  return !0;
}
function sY(J) {
  if (!J) return [];
  J = R6(J);
  let Y = {
    ctrl: "control",
    slash: "/",
    space: " ",
    spacebar: " ",
    cmd: "meta",
    esc: "escape",
    up: "arrow-up",
    down: "arrow-down",
    left: "arrow-left",
    right: "arrow-right",
    period: ".",
    comma: ",",
    equal: "=",
    minus: "-",
    underscore: "_",
  };
  return (
    (Y[J] = J),
    Object.keys(Y)
      .map((Z) => {
        if (Y[Z] === J) return Z;
      })
      .filter((Z) => Z)
  );
}
O("model", (J, { modifiers: Y, expression: Z }, { effect: $, cleanup: W }) => {
  let X = J;
  if (Y.includes("parent")) X = J.parentNode;
  let Q = I(X, Z),
    G;
  if (typeof Z === "string") G = I(X, `${Z} = __placeholder`);
  else if (typeof Z === "function" && typeof Z() === "string")
    G = I(X, `${Z()} = __placeholder`);
  else G = () => {};
  let B = () => {
      let K;
      return Q((q) => (K = q)), rY(K) ? K.get() : K;
    },
    U = (K) => {
      let q;
      if ((Q((j) => (q = j)), rY(q))) q.set(K);
      else G(() => {}, { scope: { __placeholder: K } });
    };
  if (typeof Z === "string" && J.type === "radio")
    C(() => {
      if (!J.hasAttribute("name")) J.setAttribute("name", Z);
    });
  var _ =
    J.tagName.toLowerCase() === "select" ||
    ["checkbox", "radio"].includes(J.type) ||
    Y.includes("lazy")
      ? "change"
      : "input";
  let z = m
    ? () => {}
    : KY(J, _, Y, (K) => {
        U(sJ(J, Y, K, B()));
      });
  if (Y.includes("fill")) {
    if (
      [void 0, null, ""].includes(B()) ||
      (EY(J) && Array.isArray(B())) ||
      (J.tagName.toLowerCase() === "select" && J.multiple)
    )
      U(sJ(J, Y, { target: J }, B()));
  }
  if (!J._x_removeModelListeners) J._x_removeModelListeners = {};
  if (
    ((J._x_removeModelListeners.default = z),
    W(() => J._x_removeModelListeners.default()),
    J.form)
  ) {
    let K = KY(J.form, "reset", [], (q) => {
      TY(() => J._x_model && J._x_model.set(sJ(J, Y, { target: J }, B())));
    });
    W(() => K());
  }
  (J._x_model = {
    get() {
      return B();
    },
    set(K) {
      U(K);
    },
  }),
    (J._x_forceModelUpdate = (K) => {
      if (K === void 0 && typeof Z === "string" && Z.match(/\./)) K = "";
      (window.fromModel = !0),
        C(() => kZ(J, "value", K)),
        delete window.fromModel;
    }),
    $(() => {
      let K = B();
      if (Y.includes("unintrusive") && document.activeElement.isSameNode(J))
        return;
      J._x_forceModelUpdate(K);
    });
});
function sJ(J, Y, Z, $) {
  return C(() => {
    if (Z instanceof CustomEvent && Z.detail !== void 0)
      return Z.detail !== null && Z.detail !== void 0
        ? Z.detail
        : Z.target.value;
    else if (EY(J))
      if (Array.isArray($)) {
        let W = null;
        if (Y.includes("number")) W = rJ(Z.target.value);
        else if (Y.includes("boolean")) W = IJ(Z.target.value);
        else W = Z.target.value;
        return Z.target.checked
          ? $.includes(W)
            ? $
            : $.concat([W])
          : $.filter((X) => !j6(X, W));
      } else return Z.target.checked;
    else if (J.tagName.toLowerCase() === "select" && J.multiple) {
      if (Y.includes("number"))
        return Array.from(Z.target.selectedOptions).map((W) => {
          let X = W.value || W.text;
          return rJ(X);
        });
      else if (Y.includes("boolean"))
        return Array.from(Z.target.selectedOptions).map((W) => {
          let X = W.value || W.text;
          return IJ(X);
        });
      return Array.from(Z.target.selectedOptions).map((W) => {
        return W.value || W.text;
      });
    } else {
      let W;
      if (hZ(J))
        if (Z.target.checked) W = Z.target.value;
        else W = $;
      else W = Z.target.value;
      if (Y.includes("number")) return rJ(W);
      else if (Y.includes("boolean")) return IJ(W);
      else if (Y.includes("trim")) return W.trim();
      else return W;
    }
  });
}
function rJ(J) {
  let Y = J ? parseFloat(J) : null;
  return O6(Y) ? Y : J;
}
function j6(J, Y) {
  return J == Y;
}
function O6(J) {
  return !Array.isArray(J) && !isNaN(J);
}
function rY(J) {
  return (
    J !== null &&
    typeof J === "object" &&
    typeof J.get === "function" &&
    typeof J.set === "function"
  );
}
O("cloak", (J) =>
  queueMicrotask(() => C(() => J.removeAttribute(XJ("cloak")))),
);
EZ(() => `[${XJ("init")}]`);
O(
  "init",
  p((J, { expression: Y }, { evaluate: Z }) => {
    if (typeof Y === "string") return !!Y.trim() && Z(Y, {}, !1);
    return Z(Y, {}, !1);
  }),
);
O("text", (J, { expression: Y }, { effect: Z, evaluateLater: $ }) => {
  let W = $(Y);
  Z(() => {
    W((X) => {
      C(() => {
        J.textContent = X;
      });
    });
  });
});
O("html", (J, { expression: Y }, { effect: Z, evaluateLater: $ }) => {
  let W = $(Y);
  Z(() => {
    W((X) => {
      C(() => {
        (J.innerHTML = X), (J._x_ignoreSelf = !0), u(J), delete J._x_ignoreSelf;
      });
    });
  });
});
CY(AZ(":", PZ(XJ("bind:"))));
var B0 = (
  J,
  { value: Y, modifiers: Z, expression: $, original: W },
  { effect: X, cleanup: Q },
) => {
  if (!Y) {
    let B = {};
    E1(B),
      I(J, $)(
        (_) => {
          gZ(J, _, W);
        },
        { scope: B },
      );
    return;
  }
  if (Y === "key") return T6(J, $);
  if (
    J._x_inlineBindings &&
    J._x_inlineBindings[Y] &&
    J._x_inlineBindings[Y].extract
  )
    return;
  let G = I(J, $);
  X(() =>
    G((B) => {
      if (B === void 0 && typeof $ === "string" && $.match(/\./)) B = "";
      C(() => kZ(J, Y, B, Z));
    }),
  ),
    Q(() => {
      J._x_undoAddedClasses && J._x_undoAddedClasses(),
        J._x_undoAddedStyles && J._x_undoAddedStyles();
    });
};
B0.inline = (J, { value: Y, modifiers: Z, expression: $ }) => {
  if (!Y) return;
  if (!J._x_inlineBindings) J._x_inlineBindings = {};
  J._x_inlineBindings[Y] = { expression: $, extract: !1 };
};
O("bind", B0);
function T6(J, Y) {
  J._x_keyExpression = Y;
}
VZ(() => `[${XJ("data")}]`);
O("data", (J, { expression: Y }, { cleanup: Z }) => {
  if (V6(J)) return;
  Y = Y === "" ? "{}" : Y;
  let $ = {};
  YY($, J);
  let W = {};
  I1(W, $);
  let X = r(J, Y, { scope: W });
  if (X === void 0 || X === !0) X = {};
  YY(X, J);
  let Q = $J(X);
  GZ(Q);
  let G = qJ(J, Q);
  Q.init && r(J, Q.init),
    Z(() => {
      Q.destroy && r(J, Q.destroy), G();
    });
});
yJ((J, Y) => {
  if (J._x_dataStack)
    (Y._x_dataStack = J._x_dataStack),
      Y.setAttribute("data-has-alpine-state", !0);
});
function V6(J) {
  if (!m) return !1;
  if (BY) return !0;
  return J.hasAttribute("data-has-alpine-state");
}
O("show", (J, { modifiers: Y, expression: Z }, { effect: $ }) => {
  let W = I(J, Z);
  if (!J._x_doHide)
    J._x_doHide = () => {
      C(() => {
        J.style.setProperty(
          "display",
          "none",
          Y.includes("important") ? "important" : void 0,
        );
      });
    };
  if (!J._x_doShow)
    J._x_doShow = () => {
      C(() => {
        if (J.style.length === 1 && J.style.display === "none")
          J.removeAttribute("style");
        else J.style.removeProperty("display");
      });
    };
  let X = () => {
      J._x_doHide(), (J._x_isShown = !1);
    },
    Q = () => {
      J._x_doShow(), (J._x_isShown = !0);
    },
    G = () => setTimeout(Q),
    B = QY(
      (z) => (z ? Q() : X()),
      (z) => {
        if (typeof J._x_toggleAndCascadeWithTransitions === "function")
          J._x_toggleAndCascadeWithTransitions(J, z, Q, X);
        else z ? G() : X();
      },
    ),
    U,
    _ = !0;
  $(() =>
    W((z) => {
      if (!_ && z === U) return;
      if (Y.includes("immediate")) z ? G() : X();
      B(z), (U = z), (_ = !1);
    }),
  );
});
O("for", (J, { expression: Y }, { effect: Z, cleanup: $ }) => {
  let W = N6(Y),
    X = I(J, W.items),
    Q = I(J, J._x_keyExpression || "index");
  (J._x_prevKeys = []),
    (J._x_lookup = {}),
    Z(() => E6(J, W, X, Q)),
    $(() => {
      Object.values(J._x_lookup).forEach((G) =>
        C(() => {
          GJ(G), G.remove();
        }),
      ),
        delete J._x_prevKeys,
        delete J._x_lookup;
    });
});
function E6(J, Y, Z, $) {
  let W = (Q) => typeof Q === "object" && !Array.isArray(Q),
    X = J;
  Z((Q) => {
    if (I6(Q) && Q >= 0) Q = Array.from(Array(Q).keys(), (M) => M + 1);
    if (Q === void 0) Q = [];
    let { _x_lookup: G, _x_prevKeys: B } = J,
      U = [],
      _ = [];
    if (W(Q))
      Q = Object.entries(Q).map(([M, L]) => {
        let P = aY(Y, L, M, Q);
        $(
          (A) => {
            if (_.includes(A)) k("Duplicate key on x-for", J);
            _.push(A);
          },
          { scope: { index: M, ...P } },
        ),
          U.push(P);
      });
    else
      for (let M = 0; M < Q.length; M++) {
        let L = aY(Y, Q[M], M, Q);
        $(
          (P) => {
            if (_.includes(P)) k("Duplicate key on x-for", J);
            _.push(P);
          },
          { scope: { index: M, ...L } },
        ),
          U.push(L);
      }
    let z = [],
      K = [],
      q = [],
      j = [];
    for (let M = 0; M < B.length; M++) {
      let L = B[M];
      if (_.indexOf(L) === -1) q.push(L);
    }
    B = B.filter((M) => !q.includes(M));
    let T = "template";
    for (let M = 0; M < _.length; M++) {
      let L = _[M],
        P = B.indexOf(L);
      if (P === -1) B.splice(M, 0, L), z.push([T, M]);
      else if (P !== M) {
        let A = B.splice(M, 1)[0],
          H = B.splice(P - 1, 1)[0];
        B.splice(M, 0, H), B.splice(P, 0, A), K.push([A, H]);
      } else j.push(L);
      T = L;
    }
    for (let M = 0; M < q.length; M++) {
      let L = q[M];
      if (!(L in G)) continue;
      C(() => {
        GJ(G[L]), G[L].remove();
      }),
        delete G[L];
    }
    for (let M = 0; M < K.length; M++) {
      let [L, P] = K[M],
        A = G[L],
        H = G[P],
        x = document.createElement("div");
      C(() => {
        if (!H) k('x-for ":key" is undefined or invalid', X, P, G);
        H.after(x),
          A.after(H),
          H._x_currentIfEl && H.after(H._x_currentIfEl),
          x.before(A),
          A._x_currentIfEl && A.after(A._x_currentIfEl),
          x.remove();
      }),
        H._x_refreshXForScope(U[_.indexOf(P)]);
    }
    for (let M = 0; M < z.length; M++) {
      let [L, P] = z[M],
        A = L === "template" ? X : G[L];
      if (A._x_currentIfEl) A = A._x_currentIfEl;
      let H = U[P],
        x = _[P],
        D = document.importNode(X.content, !0).firstElementChild,
        w = $J(H);
      if (
        (qJ(D, w, X),
        (D._x_refreshXForScope = (E) => {
          Object.entries(E).forEach(([d, JJ]) => {
            w[d] = JJ;
          });
        }),
        C(() => {
          A.after(D), p(() => u(D))();
        }),
        typeof x === "object")
      )
        k(
          "x-for key cannot be an object, it must be a string or an integer",
          X,
        );
      G[x] = D;
    }
    for (let M = 0; M < j.length; M++)
      G[j[M]]._x_refreshXForScope(U[_.indexOf(j[M])]);
    X._x_prevKeys = _;
  });
}
function N6(J) {
  let Y = /,([^,\}\]]*)(?:,([^,\}\]]*))?$/,
    Z = /^\s*\(|\)\s*$/g,
    $ = /([\s\S]*?)\s+(?:in|of)\s+([\s\S]*)/,
    W = J.match($);
  if (!W) return;
  let X = {};
  X.items = W[2].trim();
  let Q = W[1].replace(Z, "").trim(),
    G = Q.match(Y);
  if (G) {
    if (((X.item = Q.replace(Y, "").trim()), (X.index = G[1].trim()), G[2]))
      X.collection = G[2].trim();
  } else X.item = Q;
  return X;
}
function aY(J, Y, Z, $) {
  let W = {};
  if (/^\[.*\]$/.test(J.item) && Array.isArray(Y))
    J.item
      .replace("[", "")
      .replace("]", "")
      .split(",")
      .map((Q) => Q.trim())
      .forEach((Q, G) => {
        W[Q] = Y[G];
      });
  else if (
    /^\{.*\}$/.test(J.item) &&
    !Array.isArray(Y) &&
    typeof Y === "object"
  )
    J.item
      .replace("{", "")
      .replace("}", "")
      .split(",")
      .map((Q) => Q.trim())
      .forEach((Q) => {
        W[Q] = Y[Q];
      });
  else W[J.item] = Y;
  if (J.index) W[J.index] = Z;
  if (J.collection) W[J.collection] = $;
  return W;
}
function I6(J) {
  return !Array.isArray(J) && !isNaN(J);
}
function U0() {}
U0.inline = (J, { expression: Y }, { cleanup: Z }) => {
  let $ = xJ(J);
  if (!$._x_refs) $._x_refs = {};
  ($._x_refs[Y] = J), Z(() => delete $._x_refs[Y]);
};
O("ref", U0);
O("if", (J, { expression: Y }, { effect: Z, cleanup: $ }) => {
  if (J.tagName.toLowerCase() !== "template")
    k("x-if can only be used on a <template> tag", J);
  let W = I(J, Y),
    X = () => {
      if (J._x_currentIfEl) return J._x_currentIfEl;
      let G = J.content.cloneNode(!0).firstElementChild;
      return (
        qJ(G, {}, J),
        C(() => {
          J.after(G), p(() => u(G))();
        }),
        (J._x_currentIfEl = G),
        (J._x_undoIf = () => {
          C(() => {
            GJ(G), G.remove();
          }),
            delete J._x_currentIfEl;
        }),
        G
      );
    },
    Q = () => {
      if (!J._x_undoIf) return;
      J._x_undoIf(), delete J._x_undoIf;
    };
  Z(() =>
    W((G) => {
      G ? X() : Q();
    }),
  ),
    $(() => J._x_undoIf && J._x_undoIf());
});
O("id", (J, { expression: Y }, { evaluate: Z }) => {
  Z(Y).forEach((W) => L6(J, W));
});
yJ((J, Y) => {
  if (J._x_ids) Y._x_ids = J._x_ids;
});
CY(AZ("@", PZ(XJ("on:"))));
O(
  "on",
  p((J, { value: Y, modifiers: Z, expression: $ }, { cleanup: W }) => {
    let X = $ ? I(J, $) : () => {};
    if (J.tagName.toLowerCase() === "template") {
      if (!J._x_forwardEvents) J._x_forwardEvents = [];
      if (!J._x_forwardEvents.includes(Y)) J._x_forwardEvents.push(Y);
    }
    let Q = KY(J, Y, Z, (G) => {
      X(() => {}, { scope: { $event: G }, params: [G] });
    });
    W(() => Q());
  }),
);
dJ("Collapse", "collapse", "collapse");
dJ("Intersect", "intersect", "intersect");
dJ("Focus", "trap", "focus");
dJ("Mask", "mask", "mask");
function dJ(J, Y, Z) {
  O(Y, ($) =>
    k(
      `You can't use [x-${Y}] without first installing the "${J}" plugin here: https://alpinejs.dev/plugins/${Z}`,
      $,
    ),
  );
}
PJ.setEvaluator(KZ);
PJ.setReactivityEngine({ reactive: xY, effect: v1, release: u1, raw: R });
var D6 = PJ,
  FJ = D6;
function w6(J) {
  J.directive("collapse", Y),
    (Y.inline = (Z, { modifiers: $ }) => {
      if (!$.includes("min")) return;
      (Z._x_doShow = () => {}), (Z._x_doHide = () => {});
    });
  function Y(Z, { modifiers: $ }) {
    let W = _0($, "duration", 250) / 1000,
      X = _0($, "min", 0),
      Q = !$.includes("min");
    if (!Z._x_isShown) Z.style.height = `${X}px`;
    if (!Z._x_isShown && Q) Z.hidden = !0;
    if (!Z._x_isShown) Z.style.overflow = "hidden";
    let G = (U, _) => {
        let z = J.setStyles(U, _);
        return _.height ? () => {} : z;
      },
      B = {
        transitionProperty: "height",
        transitionDuration: `${W}s`,
        transitionTimingFunction: "cubic-bezier(0.4, 0.0, 0.2, 1)",
      };
    Z._x_transition = {
      in(U = () => {}, _ = () => {}) {
        if (Q) Z.hidden = !1;
        if (Q) Z.style.display = null;
        let z = Z.getBoundingClientRect().height;
        Z.style.height = "auto";
        let K = Z.getBoundingClientRect().height;
        if (z === K) z = X;
        J.transition(
          Z,
          J.setStyles,
          { during: B, start: { height: z + "px" }, end: { height: K + "px" } },
          () => (Z._x_isShown = !0),
          () => {
            if (Math.abs(Z.getBoundingClientRect().height - K) < 1)
              Z.style.overflow = null;
          },
        );
      },
      out(U = () => {}, _ = () => {}) {
        let z = Z.getBoundingClientRect().height;
        J.transition(
          Z,
          G,
          { during: B, start: { height: z + "px" }, end: { height: X + "px" } },
          () => (Z.style.overflow = "hidden"),
          () => {
            if (((Z._x_isShown = !1), Z.style.height == `${X}px` && Q))
              (Z.style.display = "none"), (Z.hidden = !0);
          },
        );
      },
    };
  }
}
function _0(J, Y, Z) {
  if (J.indexOf(Y) === -1) return Z;
  let $ = J[J.indexOf(Y) + 1];
  if (!$) return Z;
  if (Y === "duration") {
    let W = $.match(/([0-9]+)ms/);
    if (W) return W[1];
  }
  if (Y === "min") {
    let W = $.match(/([0-9]+)px/);
    if (W) return W[1];
  }
  return $;
}
var z0 = w6;
function b6(J) {
  let Y = () => {
    let Z, $;
    try {
      $ = localStorage;
    } catch (W) {
      console.error(W),
        console.warn(
          "Alpine: $persist is using temporary storage since localStorage is unavailable.",
        );
      let X = new Map();
      $ = { getItem: X.get.bind(X), setItem: X.set.bind(X) };
    }
    return J.interceptor(
      (W, X, Q, G, B) => {
        let U = Z || `_x_${G}`,
          _ = K0(U, $) ? M0(U, $) : W;
        return (
          Q(_),
          J.effect(() => {
            let z = X();
            L0(U, z, $), Q(z);
          }),
          _
        );
      },
      (W) => {
        (W.as = (X) => {
          return (Z = X), W;
        }),
          (W.using = (X) => {
            return ($ = X), W;
          });
      },
    );
  };
  Object.defineProperty(J, "$persist", { get: () => Y() }),
    J.magic("persist", Y),
    (J.persist = (Z, { get: $, set: W }, X = localStorage) => {
      let Q = K0(Z, X) ? M0(Z, X) : $();
      W(Q),
        J.effect(() => {
          let G = $();
          L0(Z, G, X), W(G);
        });
    });
}
function K0(J, Y) {
  return Y.getItem(J) !== null;
}
function M0(J, Y) {
  let Z = Y.getItem(J, Y);
  if (Z === void 0) return;
  return JSON.parse(Z);
}
function L0(J, Y, Z) {
  Z.setItem(J, JSON.stringify(Y));
}
var q0 = b6;
window.Alpine = FJ;
FJ.plugin(z0);
FJ.plugin(q0);
FJ.start();
hY.config.globalViewTransitions = !0;
